package resource

// buvid.go — 真 buvid3 生成与激活（2026-10-06 -352 风控修复）
// 参照: xfgryujk/blivedm BUVID_INIT_URL 种cookie流程 + B站 buvid3 格式规范
// 替换 resource.go 中硬编码 "buvid3=hi"（风控必杀假指纹）

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// genBuvid3 生成 B 站格式的 buvid3: {uuid4}${infoc_version}
// 真实样例: "C9B...F83infoc" 形如 uuid + 时间戳后缀; 实测 "xx..xx" 32hex 即可通过基本校验
func genBuvid3() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	// uuid v4 布局
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	u := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
	return u + "infoc"
}

// SeedBuvid 访问 bilibili.com 首页, 让 B 站下发正式 buvid3/buvid4 (set-cookie)
// 若已有有效 buvid3 则跳过。返回最终使用的 buvid3。
func (a *API) SeedBuvid() (string, error) {
	if a.buvid3 != "" {
		return a.buvid3, nil // 已激活
	}
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "https://www.bilibili.com/", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	got := ""
	for _, c := range resp.Cookies() {
		if c.Name == "buvid3" && c.Value != "" {
			got = c.Value
			break
		}
	}
	if got == "" {
		// 首页未下发(少见), 用本地生成的兜底
		got = genBuvid3()
	}
	a.buvid3 = got
	return got, nil
}

// BuvidCookieHeader 返回 buvid3=... 形式的 cookie 片段, 供各 client SetHeader 使用
func (a *API) BuvidCookieHeader() string {
	if a.buvid3 == "" {
		v, err := a.SeedBuvid()
		if err != nil {
			v = genBuvid3() // 网络失败也要有合法格式
			a.buvid3 = v
		}
	}
	return "buvid3=" + a.buvid3 + "; b_nut=" + fmt.Sprintf("%d", time.Now().Unix())
}

// IsBuvidSeeded 供调用方判断
func (a *API) IsBuvidSeeded() bool { return strings.TrimSpace(a.buvid3) != "" }
