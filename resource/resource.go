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
	UserAgentValue = "curl/7.68.0"
	AcceptKey      = "Accept"
	AcceptValue    = "application/json, text/plain, */*"
	CookieKey      = "Cookie"
	CookieValue    = "say=hi"
)

var (
	// 直播
	liveAPIClient = newClient().SetBaseURL(LiveAPIURL)
	// 用户信息
	apiClient = newClient().SetBaseURL(APIURL)
	// 动态
	vcApiClient = newClient().SetBaseURL(VcAPIURL)
	// 投稿
	videoApiClient = newClient().SetBaseURL(APIURL)
)

func newClient() *resty.Client {
	return resty.New().SetHeader(UserAgentKey, UserAgentValue).SetHeader(CookieKey, CookieValue)
}
