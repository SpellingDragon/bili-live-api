package resource

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
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

const DefaultCookiePath = "cookie.json"
const CopiedLoginInfoCode = 1
const NearlyExpiredLoginInfoCode = 2

var lock = sync.Mutex{}

func GetCookieInfo(cookiePath string) *CookieInfo {
	cookieInfo := &CookieInfo{}
	loginInfo, err := os.ReadFile(cookiePath)
	if err != nil || len(loginInfo) == 0 {
		log.Printf("登录信息%s不存在,使用默认登录信息%s:%+v", cookiePath, DefaultCookiePath, err)
		// 如果不是，则拷贝默认的cookies.json文件到指定路径
		loginInfo, err = os.ReadFile(DefaultCookiePath) // 假设默认文件名为default_cookies.json
		if err != nil || len(loginInfo) == 0 {
			log.Printf("%s无法使用默认登录信息%s:%+v", cookiePath, DefaultCookiePath, err.Error())
			return nil
		}
		_ = jsoniter.Unmarshal(loginInfo, &cookieInfo)
		cookieInfo.Code &= CopiedLoginInfoCode
		cookieInfo.Message = "Copied"
		loginInfo, err = jsoniter.Marshal(cookieInfo)
		if err != nil {
			log.Printf("%s无法使用默认登录信息%s:%+v", cookiePath, DefaultCookiePath, err.Error())
			return nil
		}
		lock.Lock()
		defer lock.Unlock()
		err = os.WriteFile(cookiePath, loginInfo, 0644)
		if err != nil {
			log.Printf("%s无法使用默认登录信息%s:%+v", cookiePath, DefaultCookiePath, err.Error())
			return nil
		}
	} else {
		_ = jsoniter.Unmarshal(loginInfo, &cookieInfo)
	}
	// 查找SESSDATA cookie并检查其过期时间
	for _, cookie := range cookieInfo.Data.CookieInfo.Cookies {
		if cookie.Name == "SESSDATA" {
			// 获取当前时间戳
			currentTime := time.Now().Unix()
			// 计算7天前的时间戳
			sevenDaysBeforeExpires := int64(7 * 24 * 60 * 60)
			// 检查当前时间是否在expires前7天内
			if cookie.Expires-currentTime >= sevenDaysBeforeExpires {
				cookieInfo.Code &= NearlyExpiredLoginInfoCode
				cookieInfo.Message = fmt.Sprintf("expired after %d second", cookie.Expires-currentTime)
				loginInfo, err = jsoniter.Marshal(cookieInfo)
				err = os.WriteFile(cookiePath, loginInfo, 0644)
				if err != nil {
					log.Printf("%s登录信息更新异常: info:%s error:%+v", cookiePath, loginInfo, err.Error())
				}
			}
			break
		}
	}
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
	_ = jsoniter.Unmarshal(loginInfo, cookieInfo)
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
