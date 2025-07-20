package resource

import (
	"crypto/md5"
	"fmt"
	"math/rand"
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

// EncWbi2 对请求参数进行增强，添加随机数和额外参数
// 参考 Python 版本的 _enc_wbi2 函数和真实请求参数
func EncWbi2(params map[string]string) map[string]string {
	// 定义随机字符集，与Python版本保持一致
	dmRand := "ABCDEFGHIJK"

	// 生成2位随机字符串的函数
	generateRandomStr := func() string {
		// 随机选择2个不重复的字符
		indices := rand.Perm(len(dmRand))[:2]
		result := make([]byte, 2)
		for i, idx := range indices {
			result[i] = dmRand[idx]
		}
		return string(result)
	}

	// 添加增强参数
	params["dm_img_list"] = "[]"
	params["dm_img_str"] = generateRandomStr()
	params["dm_cover_img_str"] = generateRandomStr()
	params["dm_img_inter"] = `{"ds":[],"wh":[0,0,0],"of":[0,0,0]}`
	return params
}

func EncWbi(params map[string]string, imgKey, subKey string) map[string]string {
	mixinKey := getMixinKey(imgKey + subKey)
	currTime := time.Now().Unix()
	params["wts"] = fmt.Sprintf("%d", currTime) // 添加 wts 字段
	// web_location 因为没被列入参数可能炸一些接口
	if _, exists := params["web_location"]; !exists {
		params["web_location"] = "1550101"
	}
	// 按照 key 重排参数并编码，与Python版本对齐
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	query := url.Values{}
	for _, k := range keys {
		query.Add(k, fmt.Sprintf("%v", params[k]))
	}
	// 计算 w_rid
	wbiSign := fmt.Sprintf("%x", md5.Sum([]byte(query.Encode()+mixinKey)))
	params["w_rid"] = wbiSign
	return params
}

func GetWbiKeys(imgURL string, subURL string) (string, string, error) {
	imgKey := strings.Split(imgURL[strings.LastIndex(imgURL, "/")+1:], ".")[0]
	subKey := strings.Split(subURL[strings.LastIndex(subURL, "/")+1:], ".")[0]
	return imgKey, subKey, nil
}
