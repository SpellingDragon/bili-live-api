package resource

import (
	"encoding/json"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func TestUserInfo(t *testing.T) {
	// userInfo, err := UserInfo(2258389)
	// userJson, _ := json.Marshal(userInfo)
	// println(string(userJson))
	// if err != nil {
	// 	t.Error(err)
	// }
	followerInfo, err := FollowerInfo(2258389)
	followerJson, _ := json.Marshal(followerInfo)
	println(string(followerJson))
	if err != nil {
		t.Error(err)
	}
	if followerInfo.Message == "" {
		t.Errorf("获取用户信息失败, %+v", followerInfo)
	}

	video, _ := VideoInfo("BV18d4y1n73r")
	videoJson, _ := jsoniter.Marshal(video)
	println(string(videoJson))
	tags, _ := VideoTags(video.Data)
	println(tags.Data[1].TagName)
	tagJson, _ := jsoniter.MarshalToString(tags)
	println(tagJson)
}
