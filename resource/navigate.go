package resource

import (
	"time"
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
		return nil, err
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
