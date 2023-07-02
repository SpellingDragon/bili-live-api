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
	IsPortrait           bool        `yaml:"is_portrait"`
	IsAnchor             int         `yaml:"is_anchor"`
	PkId                 int         `yaml:"pk_id"`
	RoomId               int         `yaml:"room_id"`
	ShortId              int         `yaml:"short_id"`
	LiveStatus           int         `yaml:"live_status"`
	ParentAreaId         int         `yaml:"parent_area_id"`
	NewPendants          NewPendants `yaml:"new_pendants"`
	AllowUploadCoverTime int         `yaml:"allow_upload_cover_time"`
	Description          string      `yaml:"description"`
	Background           string      `yaml:"background"`
	Title                string      `yaml:"title"`
	Tags                 string      `yaml:"tags"`
	RoomSilentSecond     int         `yaml:"room_silent_second"`
	Pendants             string      `yaml:"pendants"`
	HotWordsStatus       int         `yaml:"hot_words_status"`
	UpSession            string      `yaml:"up_session"`
	AllowChangeAreaTime  int         `yaml:"allow_change_area_time"`
	Online               int         `yaml:"online"`
	OldAreaId            int         `yaml:"old_area_id"`
	UserCover            string      `yaml:"user_cover"`
	Keyframe             string      `yaml:"keyframe"`
	IsStrictRoom         bool        `yaml:"is_strict_room"`
	RoomSilentLevel      int         `yaml:"room_silent_level"`
	AreaName             string      `yaml:"area_name"`
	PkStatus             int         `yaml:"pk_status"`
	Attention            int         `yaml:"attention"`
	Verify               string      `yaml:"verify"`
	BattleId             int         `yaml:"battle_id"`
	AreaId               int         `yaml:"area_id"`
	RoomSilentType       string      `yaml:"room_silent_type"`
	ParentAreaName       string      `yaml:"parent_area_name"`
	LiveTime             string      `yaml:"live_time"`
	StudioInfo           StudioInfo  `yaml:"studio_info"`
	Uid                  int         `yaml:"uid"`
	AreaPendants         string      `yaml:"area_pendants"`
	HotWords             []string    `yaml:"hot_words"`
}

// NewPendants
type NewPendants struct {
	Frame       Frame       `yaml:"frame"`
	Badge       interface{} `yaml:"badge"`
	MobileFrame MobileFrame `yaml:"mobile_frame"`
	MobileBadge interface{} `yaml:"mobile_badge"`
}

// Frame
type Frame struct {
	UseOldArea bool   `yaml:"use_old_area"`
	Name       string `yaml:"name"`
	Desc       string `yaml:"desc"`
	Area       int    `yaml:"area"`
	BgColor    string `yaml:"bg_color"`
	BgPic      string `yaml:"bg_pic"`
	Value      string `yaml:"value"`
	Position   int    `yaml:"position"`
	AreaOld    int    `yaml:"area_old"`
}

// MobileFrame
type MobileFrame struct {
	Name       string `yaml:"name"`
	Value      string `yaml:"value"`
	Area       int    `yaml:"area"`
	UseOldArea bool   `yaml:"use_old_area"`
	Position   int    `yaml:"position"`
	Desc       string `yaml:"desc"`
	AreaOld    int    `yaml:"area_old"`
	BgColor    string `yaml:"bg_color"`
	BgPic      string `yaml:"bg_pic"`
}

// StudioInfo
type StudioInfo struct {
	Status int `yaml:"status"`
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
		SetQueryParam("room_id", strconv.Itoa(shortID)).
		SetResult(result).
		Get("/room/v1/Room/get_info")
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
	return info.Data.RoomId, nil
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
