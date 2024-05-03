package api

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/spellingDragon/bili-live-api/dto"
	"github.com/spellingDragon/bili-live-api/log"
	"github.com/spellingDragon/bili-live-api/resource"
	"github.com/spellingDragon/bili-live-api/websocket"
)

// Live 使用 NewLive() 来初始化
type Live struct {
	Client       *websocket.Client
	ResourceAPI  *resource.API
	RoomID       int
	RoomInfo     *resource.RoomInfo
	UserInfo     *resource.UserInfo
	FollowerInfo *resource.FollowerInfo
}

const DefaultCookiePath = "cookie.json"

// NewLive 构造函数
func NewLive(roomID int) *Live {
	return &Live{
		Client:      websocket.New(DefaultCookiePath),
		ResourceAPI: resource.NewWithOptions(DefaultCookiePath, false),
		RoomID:      roomID,
	}
}

// Start 接收房间号，开始websocket心跳连接并阻塞
func (l *Live) Start() {
	for {
		if err := l.Listen(); err != nil {
			if strings.Contains(err.Error(), "获取房间号失败") {
				l.Stop()
				break
			}
			log.Warnf("连接失败, 重连中...:%v", err)
			time.Sleep(10 * time.Second)
			if err := l.Client.Connect(); err != nil {
				log.Errorf("连接websocket失败：%v", err)
				return
			}
			continue
		}
		break
	}
}

// Start 接收房间号，开始websocket心跳连接并阻塞
func (l *Live) Stop() {
	l.Client.Close()
}

func (l *Live) Listen() error {
	roomInitInfo, err := l.ResourceAPI.RoomInit(l.RoomID)
	if err != nil {
		return fmt.Errorf("获取房间号失败：%v", err)
	}

	if err := l.Client.Connect(); err != nil {
		return fmt.Errorf("连接websocket失败：%v", err)
	}

	// TODO 发送进房包,可能有顺序问题
	go l.enterRoom(roomInitInfo)

	if err := l.Client.Listening(); err != nil {
		return fmt.Errorf("监听websocket失败：%v", err)
	}
	return nil
}

// RegisterHandlers 注册不同的事件处理
// handler类型需要是定义在 websocket/handler_registration.go 中的类型，如:
// - websocket.DanmakuHandler
// - websocket.GiftHandler
// - websocket.GuardHandler
func (l *Live) RegisterHandlers(handlers ...interface{}) error {
	return l.Client.RegisterHandlers(handlers...)
}

// 发送进入房间请求
func (l *Live) enterRoom(roomInfo *resource.RoomInitResp) {
	roomInfoJson, _ := json.Marshal(roomInfo)
	log.Infof("进入房间：%s", string(roomInfoJson))
	liverInfo, err := l.ResourceAPI.GetUserInfo(roomInfo.Data.UID)
	liverInfoJson, _ := json.Marshal(liverInfo)
	log.Infof("主播信息：%s", string(liverInfoJson))
	if err != nil {
		log.Errorf("发送进入房间请求失败：%v", err)
		return
	}
	l.UserInfo = &liverInfo.Data
	getDanmu, err := l.ResourceAPI.GetDanmuInfo(roomInfo.Data.RoomID)
	if err != nil {
		log.Errorf("发送进入房间请求失败：%v", err)
		return
	}
	body, _ := jsoniter.Marshal(dto.WSEnterRoomBody{
		UID:      602310692,
		RoomID:   roomInfo.Data.RoomID, // 真实房间ID
		ProtoVer: 3,                    // 填3
		Platform: "web",
		Type:     2,
		Key:      getDanmu.Data.Token,
	})
	if err = l.Client.Write(&dto.WSPayload{
		ProtocolVersion: dto.JSON,
		Operation:       dto.RoomEnter,
		Body:            body,
	}); err != nil {
		log.Errorf("发送进入房间请求失败：%v", err)
		return
	}
}

func (l *Live) RefreshRoom() error {
	roomInfo, err := l.ResourceAPI.GetRoomInfo(l.RoomID)
	if err != nil {
		log.Errorf("刷新直播信息失败：%v", err)
		return fmt.Errorf("刷新房间信息失败：%v", err)
	}
	latestLiveTime := "0000-00-00 00:00:00"
	if l.RoomInfo != nil {
		latestLiveTime = l.RoomInfo.LiveTime
	}
	l.RoomInfo = &roomInfo.Data
	if l.RoomInfo.LiveTime == "0000-00-00 00:00:00" {
		l.RoomInfo.LiveTime = latestLiveTime
	}
	liverInfo, err := l.ResourceAPI.GetUserInfo(roomInfo.Data.UID)
	if err != nil {
		log.Errorf("刷新主播信息失败：%v", err)
		return fmt.Errorf("刷新主播信息失败：%v", err)
	}
	liverInfoJson, err := json.Marshal(liverInfo)
	if err != nil {
		log.Errorf("刷新主播信息失败：%v", err)
		return fmt.Errorf("刷新主播信息失败：%v", err)
	}
	log.Infof("主播信息：%s", string(liverInfoJson))
	l.UserInfo = &liverInfo.Data
	followerInfo, err := l.ResourceAPI.GetFollowerInfo(roomInfo.Data.UID)
	if err != nil {
		return fmt.Errorf("刷新主播粉丝数失败：%v", err)
	}
	l.FollowerInfo = &followerInfo.Data
	return nil
}

func (l *Live) GetStreamURL(qn int) string {
	playURL, err := l.ResourceAPI.GetPlayURL(l.RoomID, qn)
	if err != nil {
		log.Errorf("获取直播推流链接失败：%v", err)
		return ""
	}
	if len(playURL.Data.Durl) == 0 {
		log.Errorf("获取直播推流链接失败：%v", playURL)
		return ""
	}
	return strings.ReplaceAll(playURL.Data.Durl[0].Url, "\\u0026", "&")
}
