package api

import "testing"

func TestLive(t *testing.T) {
	live := Live{
		RoomID: 24441860,
	}
	err := live.RefreshRoom()
	if err != nil {
		t.Error(err)
	}
	streamURL := live.GetStreamURL(150)
	println(streamURL)
}
