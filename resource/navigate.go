package resource

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/spellingDragon/bili-live-api/log"
)

type NavResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		WbiImg WbiImg `json:"wbi_img"`
	} `json:"data"`
}

type WbiImg struct {
	ImgUrl string `json:"img_url"`
	SubUrl string `json:"sub_url"`
}

// RenderDataResp 用户动态页面渲染数据响应
type RenderDataResp struct {
	AccessID string `json:"access_id"`
}

// JWTPayload JWT载荷结构
type JWTPayload struct {
	Iat int64 `json:"iat"` // 创建时间
	TTL int64 `json:"ttl"` // 生存时间
}

// UserRenderData 用户渲染数据缓存
type UserRenderData struct {
	accessIDs     map[int64]string
	lastTimestamp map[int64]int64
	mutex         sync.RWMutex
}

var (
	// 全局用户渲染数据缓存实例
	userRenderDataCache = &UserRenderData{
		accessIDs:     make(map[int64]string),
		lastTimestamp: make(map[int64]int64),
	}
	// 渲染数据正则表达式
	renderDataPattern = regexp.MustCompile(`<script id="__RENDER_DATA__" type="application/json">(.*?)</script>`)
)

// SetNavCacheTTL 设置Nav缓存的TTL
func (a *API) SetNavCacheTTL(ttl time.Duration) {
	a.navMutex.Lock()
	defer a.navMutex.Unlock()
	a.navCacheTTL = ttl
}

// ClearNavCache 清除Nav缓存
func (a *API) ClearNavCache() {
	a.navMutex.Lock()
	defer a.navMutex.Unlock()
	a.navCache = nil
	a.navCacheTime = time.Time{}
}

func (a *API) Nav() (*NavResp, error) {
	// 先尝试读取缓存
	a.navMutex.RLock()
	if a.navCache != nil && time.Since(a.navCacheTime) < a.navCacheTTL {
		cachedNav := a.navCache
		a.navMutex.RUnlock()
		return cachedNav, nil
	}
	a.navMutex.RUnlock()

	// 缓存过期或不存在，重新请求
	a.navMutex.Lock()
	defer a.navMutex.Unlock()

	// 双重检查，防止并发情况下重复请求
	if a.navCache != nil && time.Since(a.navCacheTime) < a.navCacheTTL {
		return a.navCache, nil
	}

	navInfo := &NavResp{}
	_, err := a.CommonAPIClient.R().
		SetResult(navInfo).
		Get("/x/web-interface/nav")
	if err != nil {
		log.Errorf("获取导航信息失败:%+v", err)
		time.Sleep(time.Second)
		_, err = a.CommonAPIClient.R().
			SetResult(navInfo).
			Get("/x/web-interface/nav")
		if err != nil {
			log.Errorf("重新获取导航信息失败，继续尝试使用缓存:%+v", err)
			return a.navCache, nil
		}
	}

	// 更新缓存
	a.navCache = navInfo
	a.navCacheTime = time.Now()

	return navInfo, nil
}

func (a *API) GetWRID(params map[string]interface{}) (int64, string) {
	navInfo, _ := a.Nav() // 现在会使用缓存
	mainKey, subKey, _ := GetWbiKeys(navInfo.Data.WbiImg.ImgUrl, navInfo.Data.WbiImg.SubUrl)
	return EncWbi(params, mainKey, subKey)
}

// GetUserDynamicRenderData 获取用户动态页面加载静态渲染数据
// 获取部分接口需要的 w_webid 关键参数
func (a *API) GetUserDynamicRenderData(uid int64) string {
	// 检查缓存
	userRenderDataCache.mutex.RLock()
	accessID, exists := userRenderDataCache.accessIDs[uid]
	lastTime, timeExists := userRenderDataCache.lastTimestamp[uid]
	userRenderDataCache.mutex.RUnlock()

	if exists && timeExists && lastTime > time.Now().Unix() {
		return accessID
	}

	// 构建动态页面URL
	dynamicURL := fmt.Sprintf("https://space.bilibili.com/%d/dynamic", uid)

	// 发起HTTP请求
	resp, err := a.CommonAPIClient.R().Get(dynamicURL)
	if err != nil {
		log.Errorf("获取用户动态页面失败: %+v", err)
		return ""
	}

	log.Debugf("获取用户动态页面响应: %s", resp.String())
	if resp.StatusCode() != 200 {
		log.Errorf("获取用户动态页面失败，状态码: %d", resp.StatusCode())
		return ""
	}

	// 提取渲染数据
	responseText := resp.String()
	matches := renderDataPattern.FindStringSubmatch(responseText)
	if len(matches) < 2 {
		// 有的时候无需 w_webid
		return ""
	}

	// 解码URL编码的JSON数据
	scriptRenderData, err := url.QueryUnescape(matches[1])
	log.Debugf("解码后的渲染数据: %s", scriptRenderData)
	if err != nil {
		log.Errorf("URL解码失败: %+v", err)
		return ""
	}

	// 解析JSON获取access_id
	var renderData RenderDataResp
	err = json.Unmarshal([]byte(scriptRenderData), &renderData)
	if err != nil {
		log.Errorf("序列化用户动态页渲染数据异常: %+v", err)
		return ""
	}

	// 解析JWT获取过期时间（简化版本，不验证签名）
	payload, err := parseJWTPayload(renderData.AccessID)
	if err != nil {
		log.Errorf("解析JWT失败: %+v", err)
		return ""
	}

	// 更新缓存
	userRenderDataCache.mutex.Lock()
	userRenderDataCache.accessIDs[uid] = renderData.AccessID
	userRenderDataCache.lastTimestamp[uid] = payload.Iat + payload.TTL
	userRenderDataCache.mutex.Unlock()

	return renderData.AccessID
}

// parseJWTPayload 解析JWT载荷（不验证签名）
func parseJWTPayload(token string) (*JWTPayload, error) {
	// JWT格式: header.payload.signature
	parts := regexp.MustCompile(`\.`).Split(token, -1)
	if len(parts) != 3 {
		return nil, fmt.Errorf("无效的JWT格式")
	}

	// Base64解码载荷部分
	payloadBytes, err := base64Decode(parts[1])
	if err != nil {
		return nil, fmt.Errorf("Base64解码失败: %v", err)
	}

	// 解析JSON
	var payload JWTPayload
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v", err)
	}

	return &payload, nil
}

// base64Decode Base64解码（处理URL安全的Base64）
func base64Decode(data string) ([]byte, error) {
	// 替换URL安全字符
	data = strings.ReplaceAll(data, "-", "+")
	data = strings.ReplaceAll(data, "_", "/")

	// 添加必要的填充
	switch len(data) % 4 {
	case 2:
		data += "=="
	case 3:
		data += "="
	}

	// 使用标准库解码
	return base64.StdEncoding.DecodeString(data)
}
