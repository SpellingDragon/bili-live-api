package resource

import (
	"encoding/json"
	"testing"
)

func TestUserInfo(t *testing.T) {
	userInfo, err := UserInfo(2075179777)
	userJson, _ := json.Marshal(userInfo)
	println(string(userJson))

	if err != nil {
		t.Error(err)
	}
	if userInfo.Message == "" {
		t.Errorf("获取用户信息失败, %+v", userInfo)
	}

	video, err := VideoInfo("BV1tW4y1t7D6")
	if err != nil {
		t.Error(err)
	}
	videoJson, _ := json.Marshal(video)
	println(string(videoJson))
}
