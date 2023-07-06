package resource

import "github.com/go-resty/resty/v2"

const (
	// WSUrl B站直播websocket接入地址
	WSUrl = "wss://broadcastlv.chat.bilibili.com/sub"
	// LiveAPIURL B站直播API地址
	LiveAPIURL = "https://api.live.bilibili.com"
	// APIURL B站API地址
	APIURL         = "https://api.bilibili.com"
	VcAPIURL       = "https://api.vc.bilibili.com"
	UserAgentKey   = "User-Agent"
	UserAgentValue = "mengzhongshenjun/0.0.1-beta"
	AcceptKey      = "Accept"
	AcceptValue    = "application/json, text/plain, */*"
	CookieKey      = "Cookie"
	CookieValue    = "buvid3=hi"
)

var (
	// 通用
	liveAPIClient = newClient().SetDebug(false).SetBaseURL(LiveAPIURL)
	// 用户信息
	apiClient = newClient().SetDebug(false).SetBaseURL(APIURL)
	// 动态
	vcApiClient = newClient().SetDebug(false).SetBaseURL(VcAPIURL)
)

func newClient() *resty.Client {
	return resty.New().SetHeader(UserAgentKey, UserAgentValue).SetHeader(CookieKey, CookieValue)
}
