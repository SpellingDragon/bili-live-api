package resource

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/spellingDragon/bili-live-api/log"
)

// GetDanmuInfoRsp 进房响应
type GetDanmuInfoRsp struct {
	Code    int              `json:"code"`
	Msg     string           `json:"msg"`
	Message string           `json:"message"`
	Data    GetRoomDanmuInfo `json:"data"`
}

// GetRoomDanmuInfo
type GetRoomDanmuInfo struct {
	HostList         []HostList `json:"host_list"`
	Group            string     `json:"group"`
	BusinessId       int        `json:"business_id"`
	RefreshRowFactor float64    `json:"refresh_row_factor"`
	RefreshRate      int        `json:"refresh_rate"`
	MaxDelay         int        `json:"max_delay"`
	Token            string     `json:"token"`
}

// HostList
type HostList struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	WssPort int    `json:"wss_port"`
	WsPort  int    `json:"ws_port"`
}

// RoomInitResp 直播进房信息
type RoomInitResp struct {
	Code    int          `json:"code"`
	Msg     string       `json:"msg"`
	Message string       `json:"message"`
	Data    RoomInitInfo `json:"data"`
}

type RoomInitInfo struct {
	RoomID      int   `json:"room_id"`
	ShortID     int   `json:"short_id"`
	UID         int64 `json:"uid"`
	NeedP2P     int   `json:"need_p2p"`
	IsHidden    bool  `json:"is_hidden"`
	IsLocked    bool  `json:"is_locked"`
	HiddenTill  int   `json:"hidden_till"`
	LockTill    int   `json:"lock_till"`
	Encrypted   bool  `json:"encrypted"`
	PwdVerified bool  `json:"pwd_verified"`
	LiveTime    int64 `json:"live_time"`
	RoomShield  int   `json:"room_shield"`
	IsSp        int   `json:"is_sp"`
	SpecialType int   `json:"special_type"`
}

// RoomInfoResp 直播房间信息
type RoomInfoResp struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Message string   `json:"message"`
	Data    RoomInfo `json:"data"`
}

type RoomInfo struct {
	RoomID               int         `json:"room_id"`
	ShortID              int         `json:"short_id"`
	UID                  int64       `json:"uid"`
	IsPortrait           bool        `json:"is_portrait"`
	IsAnchor             int         `json:"is_anchor"`
	PkId                 int         `json:"pk_id"`
	LiveStatus           int         `json:"live_status"`
	ParentAreaId         int         `json:"parent_area_id"`
	NewPendants          NewPendants `json:"new_pendants"`
	AllowUploadCoverTime int         `json:"allow_upload_cover_time"`
	Description          string      `json:"description"`
	Background           string      `json:"background"`
	Title                string      `json:"title"`
	Tags                 string      `json:"tags"`
	RoomSilentSecond     int         `json:"room_silent_second"`
	Pendants             string      `json:"pendants"`
	HotWordsStatus       int         `json:"hot_words_status"`
	UpSession            string      `json:"up_session"`
	AllowChangeAreaTime  int         `json:"allow_change_area_time"`
	Online               int         `json:"online"`
	OldAreaId            int         `json:"old_area_id"`
	UserCover            string      `json:"user_cover"`
	Keyframe             string      `json:"keyframe"`
	IsStrictRoom         bool        `json:"is_strict_room"`
	RoomSilentLevel      int         `json:"room_silent_level"`
	AreaName             string      `json:"area_name"`
	PkStatus             int         `json:"pk_status"`
	Attention            int         `json:"attention"`
	Verify               string      `json:"verify"`
	BattleId             int         `json:"battle_id"`
	AreaId               int         `json:"area_id"`
	RoomSilentType       string      `json:"room_silent_type"`
	ParentAreaName       string      `json:"parent_area_name"`
	LiveTime             string      `json:"live_time"`
	StudioInfo           StudioInfo  `json:"studio_info"`
	AreaPendants         string      `json:"area_pendants"`
	HotWords             []string    `json:"hot_words"`
}

// NewPendants
type NewPendants struct {
	Frame       Frame       `json:"frame"`
	Badge       interface{} `json:"badge"`
	MobileFrame MobileFrame `json:"mobile_frame"`
	MobileBadge interface{} `json:"mobile_badge"`
}

// Frame
type Frame struct {
	UseOldArea bool   `json:"use_old_area"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Area       int    `json:"area"`
	BgColor    string `json:"bg_color"`
	BgPic      string `json:"bg_pic"`
	Value      string `json:"value"`
	Position   int    `json:"position"`
	AreaOld    int    `json:"area_old"`
}

// MobileFrame
type MobileFrame struct {
	Name       string `json:"name"`
	Value      string `json:"value"`
	Area       int    `json:"area"`
	UseOldArea bool   `json:"use_old_area"`
	Position   int    `json:"position"`
	Desc       string `json:"desc"`
	AreaOld    int    `json:"area_old"`
	BgColor    string `json:"bg_color"`
	BgPic      string `json:"bg_pic"`
}

// StudioInfo
type StudioInfo struct {
	Status int `json:"status"`
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

// RoomInit 获取直播间详细信息
func (a *API) RoomInit(shortID int) (*RoomInitResp, error) {
	params := a.GetWRID(map[string]string{
		"id": fmt.Sprintf("%d", shortID),
	})
	result := &RoomInitResp{}
	_, err := a.LiveAPIClient.R().
		EnableTrace().
		SetQueryParams(params).
		SetResult(result).
		Get("/room/v1/Room/room_init")
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetDanmuInfo 获取弹幕数据
func (a *API) GetDanmuInfo(shortID int) (*GetDanmuInfoRsp, error) {
	params := a.GetWRID(map[string]string{
		"id":           fmt.Sprintf("%d", shortID),
		"type":         "0",
		"web_location": "444.8",
	})
	result := &GetDanmuInfoRsp{}
	_, err := a.LiveAPIClient.R().
		EnableTrace().
		SetQueryParams(params).
		SetResult(result).
		Get("/xlive/web-room/v1/index/getDanmuInfo")
	log.Infof("获取弹幕数据:%+v", result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetRoomInfo 获取直播间详细信息
func (a *API) GetRoomInfo(shortID int) (*RoomInfoResp, error) {
	params := a.GetWRID(map[string]string{
		"room_id": fmt.Sprintf("%d", shortID),
	})
	result := &RoomInfoResp{}
	_, err := a.LiveAPIClient.R().
		EnableTrace().
		SetQueryParams(params).
		SetResult(result).
		Get("/room/v1/Room/get_info")
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RealRoomID 获取真实直播房间号,主要用于websocket连接
func (a *API) RealRoomID(shortID int) (int, error) {
	info, err := a.GetRoomInfo(shortID)
	if err != nil {
		return 0, err
	}
	return info.Data.RoomID, nil
}

// PlayURLRsp
type PlayURLRsp struct {
	Code    int         ` json:"code" json:"code,omitempty"`
	Message string      `json:"message" json:"message,omitempty"`
	Ttl     int         `json:"ttl" json:"ttl,omitempty"`
	Data    PlayURLData `json:"data" json:"data"`
}

// PlayURLData
type PlayURLData struct {
	CurrentQuality     int                  `json:"current_quality" json:"current_quality,omitempty"`
	AcceptQuality      []string             `json:"accept_quality" json:"accept_quality,omitempty"`
	CurrentQn          int                  `json:"current_qn" json:"current_qn,omitempty"`
	QualityDescription []QualityDescription `json:"quality_description" json:"quality_description,omitempty"`
	Durl               []Durl               `json:"durl" json:"durl,omitempty"`
}

// Durl
type Durl struct {
	Order      int    `json:"order" json:"order,omitempty"`
	StreamType int    `json:"stream_type" json:"stream_type,omitempty"`
	P2pType    int    `json:"p2p_type" json:"p_2_p_type,omitempty"`
	Url        string `json:"url" json:"url,omitempty"`
	Length     int    `json:"length" json:"length,omitempty"`
}

// QualityDescription
type QualityDescription struct {
	Qn   int    `json:"qn" json:"qn,omitempty"`
	Desc string `json:"desc" json:"desc,omitempty"`
}

// GetPlayURL 获取直播推流URL
func (a *API) GetPlayURL(shortID int, qn int) (*PlayURLRsp, error) {
	params := a.GetWRID(map[string]string{
		"cid":      fmt.Sprintf("%d", shortID),
		"otype":    "json",
		"platform": "web",
		"qn":       fmt.Sprintf("%d", qn),
	})
	result := &PlayURLRsp{}
	rsp, err := a.LiveAPIClient.R().
		EnableTrace().
		SetQueryParams(params).
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
