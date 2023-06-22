package resource

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/spellingDragon/bili-live-api/log"
)

// RoomInfoResp 直播房间信息
type RoomInfoResp struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Message string   `json:"message"`
	Data    RoomInfo `json:"data"`
}

type RoomInfo struct {
	RoomID      int   `json:"room_id"`
	ShortID     int   `json:"short_id"`
	UID         int   `json:"uid"`
	NeedP2P     int   `json:"need_p2p"`
	IsHidden    bool  `json:"is_hidden"`
	IsLocked    bool  `json:"is_locked"`
	IsPortrait  bool  `json:"is_portrait"`
	LiveStatus  int   `json:"live_status"`
	HiddenTill  int   `json:"hidden_till"`
	LockTill    int   `json:"lock_till"`
	Encrypted   bool  `json:"encrypted"`
	PwdVerified bool  `json:"pwd_verified"`
	LiveTime    int64 `json:"live_time"`
	RoomShield  int   `json:"room_shield"`
	IsSp        int   `json:"is_sp"`
	SpecialType int   `json:"special_type"`
}

// HostInfoResp 主播信息
type HostInfoResp struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Data    struct {
		RoomID int `json:"room_id"`
	} `json:"data"`
}

// GetRoomInfo 获取直播间详细信息
func GetRoomInfo(shortID int) (*RoomInfoResp, error) {
	result := &RoomInfoResp{}
	_, err := liveAPIClient.R().
		EnableTrace().
		SetQueryParam("id", strconv.Itoa(shortID)).
		SetResult(result).
		Get("/room/v1/Room/room_init")
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RealRoomID 获取真实直播房间号,主要用于websocket连接
func RealRoomID(shortID int) (int, error) {
	info, err := GetRoomInfo(shortID)
	if err != nil {
		return 0, err
	}
	return info.Data.RoomID, nil
}

// PlayURLRsp
type PlayURLRsp struct {
	Code    int         ` yaml:"code" json:"code,omitempty"`
	Message string      `yaml:"message" json:"message,omitempty"`
	Ttl     int         `yaml:"ttl" json:"ttl,omitempty"`
	Data    PlayURLData `yaml:"data" json:"data"`
}

// PlayURLData
type PlayURLData struct {
	CurrentQuality     int                  `yaml:"current_quality" json:"current_quality,omitempty"`
	AcceptQuality      []string             `yaml:"accept_quality" json:"accept_quality,omitempty"`
	CurrentQn          int                  `yaml:"current_qn" json:"current_qn,omitempty"`
	QualityDescription []QualityDescription `yaml:"quality_description" json:"quality_description,omitempty"`
	Durl               []Durl               `yaml:"durl" json:"durl,omitempty"`
}

// Durl
type Durl struct {
	Order      int    `yaml:"order" json:"order,omitempty"`
	StreamType int    `yaml:"stream_type" json:"stream_type,omitempty"`
	P2pType    int    `yaml:"p2p_type" json:"p_2_p_type,omitempty"`
	Url        string `yaml:"url" json:"url,omitempty"`
	Length     int    `yaml:"length" json:"length,omitempty"`
}

// QualityDescription
type QualityDescription struct {
	Qn   int    `yaml:"qn" json:"qn,omitempty"`
	Desc string `yaml:"desc" json:"desc,omitempty"`
}

// GetPlayURL 获取直播推流URL
func GetPlayURL(shortID int, qn int) (*PlayURLRsp, error) {
	result := &PlayURLRsp{}
	rsp, err := liveAPIClient.R().
		EnableTrace().
		SetQueryParam("cid", strconv.Itoa(shortID)).
		SetQueryParam("otype", "json").
		SetQueryParam("platform", "web").
		SetQueryParam("qn", strconv.Itoa(qn)).
		SetResult(result).
		Get("/room/v1/Room/playUrl")
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode() != http.StatusOK {
		log.Warnf("获取直播推流URL失败:%+v", rsp)
		return nil, errors.New(strconv.Itoa(rsp.StatusCode()))
	}
	return result, nil
}
