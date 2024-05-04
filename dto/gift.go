package dto

// ReceiveUserInfo
type ReceiveUserInfo struct {
	Uid   int    `yaml:"uid"`
	Uname string `yaml:"uname"`
}

type Gift struct {
	MedalInfo         MedalInfo       `yaml:"medal_info"`
	Switch            bool            `yaml:"switch"`
	Uid               int             `yaml:"uid"`
	Uname             string          `yaml:"uname"`
	FaceEffectId      int             `yaml:"face_effect_id"`
	GiftId            int             `yaml:"giftId"`
	Magnification     int             `yaml:"magnification"`
	BatchComboSend    BatchComboSend  `yaml:"batch_combo_send"`
	FloatScResourceId int             `yaml:"float_sc_resource_id"`
	ComboSend         ComboSend       `yaml:"combo_send"`
	DiscountPrice     int             `yaml:"discount_price"`
	EffectBlock       int             `yaml:"effect_block"`
	Num               int             `yaml:"num"`
	Super             int             `yaml:"super"`
	TopList           interface{}     `yaml:"top_list"`
	Rnd               string          `yaml:"rnd"`
	Action            string          `yaml:"action"`
	CoinType          string          `yaml:"coin_type"`
	CritProb          int             `yaml:"crit_prob"`
	SendMaster        interface{}     `yaml:"send_master"`
	TagImage          string          `yaml:"tag_image"`
	Price             int             `yaml:"price"`
	Gold              int             `yaml:"gold"`
	IsJoinReceiver    bool            `yaml:"is_join_receiver"`
	NameColor         string          `yaml:"name_color"`
	FaceEffectType    int             `yaml:"face_effect_type"`
	Rcost             int             `yaml:"rcost"`
	Remain            int             `yaml:"remain"`
	Silver            int             `yaml:"silver"`
	BeatId            string          `yaml:"beatId"`
	BlindGift         interface{}     `yaml:"blind_gift"`
	ComboResourcesId  int             `yaml:"combo_resources_id"`
	ComboStayTime     int             `yaml:"combo_stay_time"`
	Dmscore           int             `yaml:"dmscore"`
	GiftType          int             `yaml:"giftType"`
	SuperGiftNum      int             `yaml:"super_gift_num"`
	SvgaBlock         int             `yaml:"svga_block"`
	BagGift           interface{}     `yaml:"bag_gift"`
	GiftName          string          `yaml:"giftName"`
	OriginalGiftName  string          `yaml:"original_gift_name"`
	BroadcastId       int             `yaml:"broadcast_id"`
	Draw              int             `yaml:"draw"`
	IsFirst           bool            `yaml:"is_first"`
	BatchComboId      string          `yaml:"batch_combo_id"`
	GuardLevel        int             `yaml:"guard_level"`
	Tid               string          `yaml:"tid"`
	BizSource         string          `yaml:"biz_source"`
	Face              string          `yaml:"face"`
	ReceiveUserInfo   ReceiveUserInfo `yaml:"receive_user_info"`
	Demarcation       int             `yaml:"demarcation"`
	IsNaming          bool            `yaml:"is_naming"`
	Timestamp         int             `yaml:"timestamp"`
	ComboTotalCoin    int             `yaml:"combo_total_coin"`
	IsSpecialBatch    int             `yaml:"is_special_batch"`
	Effect            int             `yaml:"effect"`
	SuperBatchGiftNum int             `yaml:"super_batch_gift_num"`
	TotalCoin         int             `yaml:"total_coin"`
}

// MedalInfo
type MedalInfo struct {
	MedalColor       int    `yaml:"medal_color"`
	MedalColorBorder int    `yaml:"medal_color_border"`
	Special          string `yaml:"special"`
	TargetId         int    `yaml:"target_id"`
	AnchorRoomid     int    `yaml:"anchor_roomid"`
	AnchorUname      string `yaml:"anchor_uname"`
	IconId           int    `yaml:"icon_id"`
	MedalColorStart  int    `yaml:"medal_color_start"`
	MedalLevel       int    `yaml:"medal_level"`
	MedalName        string `yaml:"medal_name"`
	GuardLevel       int    `yaml:"guard_level"`
	IsLighted        int    `yaml:"is_lighted"`
	MedalColorEnd    int    `yaml:"medal_color_end"`
}

// BatchComboSend
type BatchComboSend struct {
	GiftId        int         `yaml:"gift_id"`
	GiftName      string      `yaml:"gift_name"`
	SendMaster    interface{} `yaml:"send_master"`
	BlindGift     interface{} `yaml:"blind_gift"`
	BatchComboId  string      `yaml:"batch_combo_id"`
	BatchComboNum int         `yaml:"batch_combo_num"`
	GiftNum       int         `yaml:"gift_num"`
	Uid           int         `yaml:"uid"`
	Uname         string      `yaml:"uname"`
	Action        string      `yaml:"action"`
}
