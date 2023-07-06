package websocket

import (
	"encoding/binary"
	"encoding/json"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"log"

	"github.com/spellingDragon/bili-live-api/dto"
)

type opCodeHandler func(*dto.WSPayload, *Client)

var opCodeHandlerMap = map[dto.OPCode]opCodeHandler{
	dto.Notification:      notificationHandler,
	dto.HeartbeatResponse: heartbeatResponseHandler,
	dto.RoomEnterResponse: roomEnterResponseHandler,
}

// 通知类消息，弹幕、礼物等
func notificationHandler(payload *dto.WSPayload, client *Client) {
	eType := jsoniter.Get(payload.Body, "cmd").ToString()
	var handler eventPayloadHandler
	// HACK 更新后收到的cmd会变为"DANMU_MSG:4:0:2:2:2:0", 在此特殊处理
	if strings.HasPrefix(eType, "DANMU_MSG") {
		log.Println("收到CMD消息:", eType, string(payload.Body))
		handler = eventPayloadHandlerMap[dto.EventDanmaku]
	} else {
		eventType := dto.EventType(eType)
		log.Println("收到CMD消息:", eventType, string(payload.Body))
		var ok bool
		handler, ok = eventPayloadHandlerMap[eventType]
		if !ok {
			log.Println("未知cmd类型: %s", eventType)
			handler = unknownEventHandler
		}
	}
	handler(payload, client)
}

// 心跳回应,body为人气值
func heartbeatResponseHandler(payload *dto.WSPayload, client *Client) {
	if client.DefaultEventHandlers.Popularity == nil {
		return
	}
	popularity := binary.BigEndian.Uint32(payload.Body)
	client.DefaultEventHandlers.Popularity(popularity)
	log.Println("收到心跳回应,人气值:", popularity)
}

// 进房回应，body为空
func roomEnterResponseHandler(payload *dto.WSPayload, client *Client) {
	roomJson, _ := json.Marshal(*client)
	log.Println("订阅信息：%s", roomJson)
}
