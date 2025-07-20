package resource

import (
	"encoding/json"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func TestRoomInfo(t *testing.T) {
	a := New()
	roomRsp, err := a.GetRoomInfo(24190721)
	roomJson, _ := json.Marshal(roomRsp)
	println(string(roomJson))
	if err != nil {
		t.Error(err)
	}
}

func TestUserInfo(t *testing.T) {
	a := New()
	userInfo, err := a.GetUserInfo(2075179777)
	userJson, _ := json.Marshal(userInfo)
	println(string(userJson))
	if err != nil {
		t.Error(err)
	}
}

// TestUserInfoWithRealParams 使用真实请求参数测试用户信息获取
func TestUserInfoWithRealParams(t *testing.T) {
	a := NewWithOptions("cookie.json", true) // 启用调试模式

	// 真实请求参数（来自浏览器抓包）
	realParams := map[string]string{
		"mid":              "1372936974",
		"token":            "",
		"platform":         "web",
		"web_location":     "1550101",
		"dm_img_list":      "[]",
		"dm_img_str":       "V2ViR0wgMS4wIChPcGVuR0wgRVMgMi4wIENocm9taXVtKQ",
		"dm_cover_img_str": "QU5HTEUgKEFwcGxlLCBBTkdMRSBNZXRhbCBSZW5kZXJlcjogQXBwbGUgTTMgUHJvLCBVbnNwZWNpZmllZCBWZXJzaW9uKUdvb2dsZSBJbmMuIChBcHBsZS",
		"dm_img_inter":     "%7B%22ds%22:[],%22wh%22:[3345,1815,105],%22of%22:[230,460,230]%7D",
		"w_webid":          "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzcG1faWQiOiIzMzMuMTM4NyIsImJ1dmlkIjoiRDQxMzBFQTMtQzQyMi03MDhDLTUyOTUtODhBMDlCRDBBMzYxMTYxODBpbmZvYyIsInVzZXJfYWdlbnQiOiJNb3ppbGxhLzUuMCAoTWFjaW50b3NoOyBJbnRlbCBNYWMgT1MgWCAxMF8xNV83KSBBcHBsZVdlYktpdC81MzcuMzYgKEtIVE1MLCBsaWtlIEdlY2tvKSBDaHJvbWUvMTM3LjAuMC4wIFNhZmFyaS81MzcuMzYiLCJidXZpZF9mcCI6IjNjZWVlMmU0NmM0NjA1NWRkMDFkZTk3ZTU5MTcxY2ZkIiwiYmlsaV90aWNrZXQiOiJleUpoYkdjaU9pSklVekkxTmlJc0ltdHBaQ0k2SW5Nd015SXNJblI1Y0NJNklrcFhWQ0o5LmV5SmxlSEFpT2pFM05UTXlNVEl6T1Rrc0ltbGhkQ0k2TVRjMU1qazFNekV6T1N3aWNHeDBJam90TVgwLjA3cUtMUHU4S3BuUzd2am1jYWRtdWYwSmVoQzhSa29TWXNUTWw2YUM3MVUiLCJjcmVhdGVkX2F0IjoxNzUyOTUzMjE5LCJ0dGwiOjg2NDAwLCJ1cmwiOiIvMTM3MjkzNjk3NC9keW5hbWljP3NwbV9pZF9mcm9tPTMzMy4xMzY1Lmxpc3QuY2FyZF9hdmF0YXIuY2xpY2siLCJyZXN1bHQiOjAsImlzcyI6ImdhaWEiLCJpYXQiOjE3NTI5NTMyMTl9.WlOkQ-25IIWnyczh8pJS1qtITfWAOJu0uu58H0BQk36zHtsw-HLO-jIwiI9z0Gl1GLAyhaMTVxzxriqAaLkkdwxO3uJBac6m076pmArtkYR-TENQOujSlGv7m-NZKkGpJt_Dp8j6-ezcKY4WicxgOM-O3dUUOdxQ2cDIPlsPJzsHNuPoXq1xyg6HSPmXiDtiAS6rHyugxuR6sacNP4lTCq8rEOk5pGX238amsS3kXtnOr8M_GUzr--4uD448ohtOLazdzTVhSrRxI_CHB6CY45Bb4nvkNxYZeLIQu-hkpO6sqBKP0zqDlL7yuIFi-wG-MZOxUBOmFBMtwe1zz11MOQ",
		"w_rid":            "c0755418820504bd27ac638a9179a37a",
		"wts":              "1752953220",
	}

	t.Logf("请求参数: %+v", realParams)

	// 设置真实的User-Agent和Headers
	userInfo := &UserInfoResp{}
	resp, err := a.CommonAPIClient.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36").
		SetHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7").
		SetHeader("Accept-Language", "zh-CN,zh;q=0.9").
		SetHeader("Cache-Control", "max-age=0").
		SetHeader("Sec-Ch-Ua", `"Google Chrome";v="137", "Chromium";v="137", "Not/A)Brand";v="24"`).
		SetHeader("Sec-Ch-Ua-Mobile", "?0").
		SetHeader("Sec-Ch-Ua-Platform", `"macOS"`).
		SetHeader("Sec-Fetch-Dest", "document").
		SetHeader("Sec-Fetch-Mode", "navigate").
		SetHeader("Sec-Fetch-Site", "none").
		SetHeader("Sec-Fetch-User", "?1").
		SetHeader("Upgrade-Insecure-Requests", "1").
		SetQueryParams(realParams).
		SetResult(userInfo).
		Get("/x/space/wbi/acc/info")

	if err != nil {
		t.Errorf("请求失败: %v", err)
		return
	}

	t.Logf("响应状态码: %d", resp.StatusCode())
	t.Logf("响应Headers: %+v", resp.Header())
	t.Logf("响应Body: %s", resp.String())
	t.Logf("解析后的响应: %+v", userInfo)

	if userInfo.Code != 0 {
		t.Errorf("API返回错误: code=%d, message=%s", userInfo.Code, userInfo.Message)
		return
	}

	if userInfo.Data.Mid == 0 {
		t.Error("未获取到用户信息")
		return
	}

	// 输出用户信息用于验证
	userJson, _ := json.Marshal(userInfo)
	t.Logf("真实参数请求结果: %s", string(userJson))
}

func TestFollowInfo(t *testing.T) {
	a := New()
	followerInfo, err := a.GetFollowerInfo(2075179777)
	followerJson, _ := json.Marshal(followerInfo)
	println(string(followerJson))
	if err != nil {
		t.Error(err)
	}
	if followerInfo.Message == "" {
		t.Errorf("获取用户信息失败, %+v", followerInfo)
	}
}

func TestVideoInfo(t *testing.T) {
	a := New()
	video, _ := a.VideoInfo("BV18d4y1n73r")
	// videoJson, _ := jsoniter.Marshal(video)
	// println(string(videoJson))
	tags, _ := a.VideoTags(video.Data)
	// println(tags.Data[1].TagName)
	tagJson, _ := jsoniter.MarshalToString(tags)
	println(tagJson)
}

func TestRoom(t *testing.T) {
	a := New()
	roomRsp, err := a.GetPlayURL(21013446, 150)
	roomJson, _ := json.Marshal(roomRsp)
	println(string(roomJson))
	if err != nil {
		t.Error(err)
	}
}
