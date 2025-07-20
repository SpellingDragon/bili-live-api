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
	if _, webLocationExists := params["web_location"]; !webLocationExists {
		params["web_location"] = "1550101"
	}
	// 生成更真实的 dm_img_str (模拟 WebGL 渲染器信息)
	generateDmImgStr := func() string {
		webglRenderers := []string{
			"V2ViR0wgMS4wIChPcGVuR0wgRVMgMi4wIENocm9taXVtKQ",
			"V2ViR0wgMS4wIChPcGVuR0wgRVMgMy4wIENocm9taXVtKQ",
			"V2ViR0wgMS4wIChPcGVuR0wgRVMgMi4wIEZpcmVmb3gpKQ",
		}
		return webglRenderers[rand.Intn(len(webglRenderers))]
	}

	// 生成更真实的 dm_cover_img_str (模拟 ANGLE 渲染器信息)
	generateDmCoverImgStr := func() string {
		angleRenderers := []string{
			"QU5HTEUgKEFwcGxlLCBBTkdMRSBNZXRhbCBSZW5kZXJlcjogQXBwbGUgTTMgUHJvLCBVbnNwZWNpZmllZCBWZXJzaW9uKUdvb2dsZSBJbmMuIChBcHBsZS",
			"QU5HTEUgKEludGVsLCBBTkdMRSBEaXJlY3QzRCAxMSBWUyA1XzAgUFMgNV8wLCBEMTFfMSlHb29nbGUgSW5jLiAoSW50ZWwp",
			"QU5HTEUgKE5WSURJQSwgQU5HTEUgRGlyZWN0M0QgMTEgVlMgNV8wIFBTIDVfMCwgRDExXzEpR29vZ2xlIEluYy4gKE5WSURJQS",
		}
		return angleRenderers[rand.Intn(len(angleRenderers))]
	}

	// 生成更真实的 dm_img_inter (模拟真实的屏幕和偏移参数)
	generateDmImgInter := func() string {
		// 常见的屏幕分辨率和偏移值
		screenConfigs := []struct {
			wh [3]int
			of [3]int
		}{
			{wh: [3]int{3345, 1815, 105}, of: [3]int{230, 460, 230}},
			{wh: [3]int{1920, 1080, 96}, of: [3]int{0, 0, 0}},
			{wh: [3]int{2560, 1440, 110}, of: [3]int{100, 200, 100}},
			{wh: [3]int{1366, 768, 96}, of: [3]int{50, 100, 50}},
		}
		config := screenConfigs[rand.Intn(len(screenConfigs))]
		return fmt.Sprintf(`{"ds":[],"wh":[%d,%d,%d],"of":[%d,%d,%d]}`,
			config.wh[0], config.wh[1], config.wh[2],
			config.of[0], config.of[1], config.of[2])
	}

	// 添加增强参数
	params["dm_img_list"] = "[]"
	params["dm_img_str"] = generateDmImgStr()
	params["dm_cover_img_str"] = generateDmCoverImgStr()
	params["dm_img_inter"] = generateDmImgInter()

	return params
}

func EncWbi(params map[string]string, imgKey, subKey string) map[string]string {
	mixinKey := getMixinKey(imgKey + subKey)
	currTime := time.Now().Unix()
	params["wts"] = fmt.Sprintf("%d", currTime) // 添加 wts 字段
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
	return params
}

func GetWbiKeys(imgURL string, subURL string) (string, string, error) {
	imgKey := strings.Split(imgURL[strings.LastIndex(imgURL, "/")+1:], ".")[0]
	subKey := strings.Split(subURL[strings.LastIndex(subURL, "/")+1:], ".")[0]
	return imgKey, subKey, nil
}
