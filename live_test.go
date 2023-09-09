package api

import (
	"bufio"
	"log"
	"os"
	"testing"
)

func TestLive(t *testing.T) {
	file, _ := os.Create("test.log")
	defer file.Close()
	log.SetOutput(bufio.NewWriter(file))
	live := NewLive("cookie.json", 24441860)
	// roomInit, err := resource.RoomInit(24441860)
	// if err != nil {
	// 	t.Error(err)
	// }
	live.Start()
	//
	// streamURL := live.GetStreamURL(150)
	// println(streamURL)
}
