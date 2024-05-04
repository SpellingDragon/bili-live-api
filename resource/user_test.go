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
