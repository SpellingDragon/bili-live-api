package resource

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func TestUserInfo(t *testing.T) {
	// userInfo, err := UserInfo(2075179777)
	// userJson, _ := json.Marshal(userInfo)
	// println(string(userJson))
	//
	// if err != nil {
	// 	t.Error(err)
	// }
	// if userInfo.Message == "" {
	// 	t.Errorf("获取用户信息失败, %+v", userInfo)
	// }

	video, _ := VideoInfo("BV18d4y1n73r")
	videoJson, _ := jsoniter.Marshal(video)
	println(string(videoJson))
	tags, _ := VideoTags(video.Data)
	println(tags.Data[1].TagName)
	tagJson, _ := jsoniter.MarshalToString(tags)
	println(tagJson)
}
