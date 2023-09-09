package resource

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type CookieInfo struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Ttl     int64  `json:"ttl"`
	Data    struct {
		IsNew        bool   `json:"is_new"`
		Mid          int64  `json:"mid"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
		TokenInfo    struct {
			Mid          int64  `json:"mid"`
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			ExpiresIn    int64  `json:"expires_in"`
		} `json:"token_info"`
		CookieInfo struct {
			Cookies []struct {
				Name     string `json:"name"`
				Value    string `json:"value"`
				HttpOnly int64  `json:"http_only"`
				Expires  int64  `json:"expires"`
				Secure   int64  `json:"secure"`
			} `json:"cookies"`
			Domains []string `json:"domains"`
		} `json:"cookie_info"`
		Sso []string `json:"sso"`
	} `json:"data"`
}

func GetCookieInfo(cookiePath string) *CookieInfo {
	cookieInfo := &CookieInfo{}
	loginInfo, err := os.ReadFile(cookiePath)
	if err != nil || len(loginInfo) == 0 {
		log.Printf("cookie文件不存在,请先登录:%+v", err)
		return nil
	}
	_ = json.Unmarshal(loginInfo, &cookieInfo)
	return cookieInfo
}

func GetCookie(cookiePath string) (cookie string, csrf string) {
	loginInfo, err := os.ReadFile(cookiePath)
	if err != nil || len(loginInfo) == 0 {
		log.Printf("cookie文件不存在,请先登录:%+v", err)
		return
	}
	cookieInfo := GetCookieInfo(cookiePath)
	if cookieInfo == nil {
		log.Printf("cookie文件反序列化失败")
		return
	}
	_ = json.Unmarshal(loginInfo, cookieInfo)
	for _, v := range cookieInfo.Data.CookieInfo.Cookies {
		cookie += v.Name + "=" + v.Value + ";"
		if v.Name == "bili_jct" {
			csrf = v.Value
		}
	}
	return
}

func ListHttpCookies(cookiePath string) []*http.Cookie {
	var cookies []*http.Cookie
	cookieInfo := GetCookieInfo(cookiePath)
	if cookieInfo == nil {
		return nil
	}
	for _, v := range cookieInfo.Data.CookieInfo.Cookies {
		cookies = append(cookies, &http.Cookie{
			Name:       v.Name,
			Value:      v.Value,
			Expires:    time.Unix(v.Expires, 0),
			RawExpires: strconv.Itoa(int(v.Expires)),
			Secure:     v.Secure != 0,
			HttpOnly:   v.HttpOnly != 0,
		})
	}
	return cookies
}
