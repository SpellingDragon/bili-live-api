package resource

import (

	"os"
	"strings"
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
	CookieValue    = "" // filled at runtime by BuvidCookieHeader() (was "buvid3=hi" — fake fingerprint, -352 root cause)
)

type API struct {
	CookiePath      string
	buvid3          string
	biliTicket      string
	ticketAt        time.Time
	cookie          string
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
	buvid, _ := a.SeedBuvid()
	a.LiveAPIClient = newClient(a.CookiePath).
		SetHeader(CookieKey, "buvid3="+buvid).SetBaseURL(LiveAPIURL)
	// 用户信息
	a.CommonAPIClient = newClient(a.CookiePath).
		SetHeader(RefererKey, RefererValue).
		SetBaseURL(APIURL)
	// 动态
	a.VcAPIClient = newClient(a.CookiePath).
		SetHeader(CookieKey, CookieValue).
		SetBaseURL(VcAPIURL)
	// 空间
	a.SpaceAPIClient = newClient(a.CookiePath).
		SetHeader(CookieKey, CookieValue).
		SetBaseURL(SpaceURL)
	return a
}

func NewWithOptions(path string, debug bool) *API {
	a := &API{}
	a.CookiePath = path
	a.navCacheTTL = 10 * time.Minute // 默认缓存10分钟
	// v0.9.10: read cookie file content for ticket/buvid chain
	if b, err := os.ReadFile(path); err == nil {
		a.cookie = strings.TrimSpace(string(b))
	}
	buvid, _ := a.SeedBuvid()
	_ = buvid
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
