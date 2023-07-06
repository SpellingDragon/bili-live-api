package api

import (
	"testing"

	"github.com/spellingDragon/bili-live-api/resource"
)

func TestLive(t *testing.T) {
	live := NewLive(24441860)
	roomInit, err := resource.RoomInit(24441860)
	if err != nil {
		t.Error(err)
	}
	go live.enterRoom(roomInit)
	if err := live.Client.Listening(); err != nil {
		t.Errorf("监听websocket失败：%v", err)
	}
	//
	// streamURL := live.GetStreamURL(150)
	// println(streamURL)
}
