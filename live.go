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

// NewLive 构造函数
func NewLive(roomID int) *Live {
	return &Live{
		Client:      websocket.New(resource.DefaultCookiePath),
		ResourceAPI: resource.NewWithOptions(resource.DefaultCookiePath, false),
		RoomID:      roomID,
	}
}

// NewLiveWithPath 构造函数
func NewLiveWithPath(roomID int, path string, debug bool) *Live {
	return &Live{
		Client:      websocket.New(path),
		ResourceAPI: resource.NewWithOptions(path, debug),
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
			// 移除这里的硬编码连接，因为连接逻辑已经移到Listen方法中
			continue
		}
		break
	}
}

// Stop 停止websocket连接
func (l *Live) Stop() {
	l.Client.Close()
}

func (l *Live) Listen() error {
	roomInitInfo, err := l.ResourceAPI.RoomInit(l.RoomID)
	if err != nil {
		return fmt.Errorf("获取房间号失败：%v", err)
	}

	// 获取弹幕服务器配置
	getDanmu, err := l.ResourceAPI.GetDanmuInfo(roomInitInfo.Data.RoomID)
	if err != nil {
		return fmt.Errorf("获取弹幕服务器配置失败：%v", err)
	}
	
	// 动态选择WebSocket服务器
	hostList := getDanmu.Data.HostList
	if len(hostList) == 0 {
		return fmt.Errorf("没有可用的WebSocket服务器")
	}
	
	// 尝试连接可用的主机（倒序尝试，与Python版本保持一致）
	var wsUrl string
	var connectErr error
	for i := len(hostList) - 1; i >= 0; i-- {
		host := hostList[i]
		wsUrl = fmt.Sprintf("wss://%s:%d/sub", host.Host, host.WssPort)
		log.Infof("正在尝试连接主机：%s", wsUrl)
		
		connectErr = l.Client.Connect(wsUrl)
		if connectErr == nil {
			log.Infof("成功连接到主机：%s", wsUrl)
			break
		}
		log.Warnf("连接主机失败：%s, 错误：%v", wsUrl, connectErr)
	}
	
	if connectErr != nil {
		return fmt.Errorf("所有主机连接失败，最后错误：%v", connectErr)
	}

	// 发送进房包
	go l.enterRoom(roomInitInfo, getDanmu)

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
func (l *Live) enterRoom(roomInfo *resource.RoomInitResp, danmuInfo *resource.GetDanmuInfoRsp) {
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
	
	// 从cookie中获取用户ID
	uid, err := resource.GetUserIDFromCookie(l.Client.CookiePath)
	if err != nil {
		log.Warnf("从cookie获取用户ID失败，使用默认值0: %v", err)
		uid = 0
	}
	log.Debugf("使用用户ID: %d", uid)
	
	// 使用传入的弹幕信息，避免重复获取
	body, _ := jsoniter.Marshal(dto.WSEnterRoomBody{
		UID:      uid,                    // 使用从cookie解析的用户ID
		RoomID:   roomInfo.Data.RoomID,   // 真实房间ID
		ProtoVer: 3,                      // 填3
		Platform: "web",
		Type:     2,
		Key:      danmuInfo.Data.Token,
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
