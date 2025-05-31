package resource

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

var mixinKeyEncTab = []int{
	46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35, 27, 43, 5, 49,
	33, 9, 42, 19, 29, 28, 14, 39, 12, 38, 41, 13, 37, 48, 7, 16, 24, 55, 40,
	61, 26, 17, 0, 1, 60, 51, 30, 4, 22, 25, 54, 21, 56, 59, 6, 63, 57, 62, 11,
	36, 20, 34, 44, 52,
}

func getMixinKey(orig string) string {
	var s string
	for _, i := range mixinKeyEncTab {
		s += orig[i : i+1]
	}
	return s[:32]
}

func EncWbi(params map[string]interface{}, imgKey, subKey string) (int64, string) {
	mixinKey := getMixinKey(imgKey + subKey)
	currTime := time.Now().Unix()
	params["wts"] = currTime // 添加 wts 字段
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys) // 按照 key 重排参数
	// 过滤 value 中的 "!'()*" 字符
	for k, v := range params {
		params[k] = strings.Map(func(r rune) rune {
			if strings.IndexRune("!'()*", r) < 0 {
				return r
			}
			return -1
		}, fmt.Sprintf("%v", v))
	}
	query := url.Values{}
	for _, k := range keys {
		query.Add(k, fmt.Sprintf("%v", params[k]))
	}
	// 计算 w_rid
	wbiSign := fmt.Sprintf("%x", md5.Sum([]byte(query.Encode()+mixinKey)))
	params["w_rid"] = wbiSign
	return currTime, wbiSign
}

func GetWbiKeys(imgURL string, subURL string) (string, string, error) {
	imgKey := strings.Split(imgURL[strings.LastIndex(imgURL, "/")+1:], ".")[0]
	subKey := strings.Split(subURL[strings.LastIndex(subURL, "/")+1:], ".")[0]
	return imgKey, subKey, nil
}
