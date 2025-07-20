package resource

import (
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	// WSUrl B站直播websocket接入地址
	WSUrl = "wss://broadcastlv.chat.bilibili.com/sub"
	// LiveAPIURL B站直播API地址
	LiveAPIURL = "https://api.live.bilibili.com"
	// APIURL B站API地址
	APIURL         = "https://api.bilibili.com"
	SpaceURL       = "https://space.bilibili.com"
	VcAPIURL       = "https://api.vc.bilibili.com"
	UserAgentKey   = "User-Agent"
	UserAgentValue = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0"
	AcceptKey      = "Accept"
	AcceptValue    = "application/json, text/plain, */*"
	RefererKey     = "Referer"
	RefererValue   = "https://www.bilibili.com"
	CookieKey      = "Cookie"
	CookieValue    = "buvid3=hi"
)

type API struct {
	CookiePath      string
	LiveAPIClient   *resty.Client
	CommonAPIClient *resty.Client
	SpaceAPIClient  *resty.Client
	VcAPIClient     *resty.Client
	// Nav缓存相关字段
	navCache     *NavResp
	navCacheTime time.Time
	navCacheTTL  time.Duration
	navMutex     sync.RWMutex
}

func New() *API {
	a := &API{}
	a.CookiePath = "cookie.json"
	a.navCacheTTL = 10 * time.Minute // 默认缓存10分钟
	// 通用
	a.LiveAPIClient = newClient(a.CookiePath).SetDebug(false).SetBaseURL(LiveAPIURL)
	// 用户信息
	a.CommonAPIClient = newClient(a.CookiePath).SetDebug(false).SetBaseURL(APIURL)
	// 动态
	a.VcAPIClient = newClient(a.CookiePath).SetDebug(false).SetBaseURL(VcAPIURL)
	// 空间
	a.SpaceAPIClient = newClient(a.CookiePath).SetDebug(false).SetBaseURL(SpaceURL)
	return a
}

func NewWithOptions(path string, debug bool) *API {
	a := &API{}
	a.CookiePath = path
	a.navCacheTTL = 10 * time.Minute // 默认缓存10分钟
	// 通用
	a.LiveAPIClient = newClient(a.CookiePath).
		SetHeader(CookieKey, CookieValue).
		SetDebug(debug).SetBaseURL(LiveAPIURL)
	// 用户信息
	a.CommonAPIClient = newClient(a.CookiePath).
		SetHeader(RefererKey, RefererValue).
		SetDebug(debug).SetBaseURL(APIURL)
	// 动态
	a.VcAPIClient = newClient(a.CookiePath).SetDebug(debug).SetBaseURL(VcAPIURL)
	// 空间
	a.SpaceAPIClient = newClient(a.CookiePath).
		SetHeader(RefererKey, RefererValue).
		SetDebug(debug).SetBaseURL(SpaceURL)
	return a
}

func newClient(cookiePath string) *resty.Client {
	return resty.New().
		SetHeader(UserAgentKey, UserAgentValue).
		SetHeader(AcceptKey, AcceptValue).
		SetCookies(ListHttpCookies(cookiePath))
}
