package resource

import (
	"errors"
	"strconv"
)

type UserInfoResp struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	TTL     int      `json:"ttl"`
	Data    UserInfo `json:"data"`
}

type UserInfo struct {
	Mid       int    `json:"mid"`
	Name      string `json:"name"`
	Sex       string `json:"sex"`
	Face      string `json:"face"`
	FaceNft   int    `json:"face_nft"`
	Sign      string `json:"sign"`
	Rank      int    `json:"rank"`
	Level     int    `json:"level"`
	Jointime  int    `json:"jointime"`
	Moral     int    `json:"moral"`
	Silence   int    `json:"silence"`
	Coins     int    `json:"coins"`
	FansBadge bool   `json:"fans_badge"`
	FansMedal struct {
		Show  bool        `json:"show"`
		Wear  bool        `json:"wear"`
		Medal interface{} `json:"medal"`
	} `json:"fans_medal"`
	Official struct {
		Role  int    `json:"role"`
		Title string `json:"title"`
		Desc  string `json:"desc"`
		Type  int    `json:"type"`
	} `json:"official"`
	Vip struct {
		Type       int   `json:"type"`
		Status     int   `json:"status"`
		DueDate    int64 `json:"due_date"`
		VipPayType int   `json:"vip_pay_type"`
		ThemeType  int   `json:"theme_type"`
		Label      struct {
			Path        string `json:"path"`
			Text        string `json:"text"`
			LabelTheme  string `json:"label_theme"`
			TextColor   string `json:"text_color"`
			BgStyle     int    `json:"bg_style"`
			BgColor     string `json:"bg_color"`
			BorderColor string `json:"border_color"`
		} `json:"label"`
		AvatarSubscript    int    `json:"avatar_subscript"`
		NicknameColor      string `json:"nickname_color"`
		Role               int    `json:"role"`
		AvatarSubscriptURL string `json:"avatar_subscript_url"`
	} `json:"vip"`
	Pendant struct {
		Pid               int    `json:"pid"`
		Name              string `json:"name"`
		Image             string `json:"image"`
		Expire            int    `json:"expire"`
		ImageEnhance      string `json:"image_enhance"`
		ImageEnhanceFrame string `json:"image_enhance_frame"`
	} `json:"pendant"`
	Nameplate struct {
		Nid        int    `json:"nid"`
		Name       string `json:"name"`
		Image      string `json:"image"`
		ImageSmall string `json:"image_small"`
		Level      string `json:"level"`
		Condition  string `json:"condition"`
	} `json:"nameplate"`
	UserHonourInfo struct {
		Mid    int           `json:"mid"`
		Colour interface{}   `json:"colour"`
		Tags   []interface{} `json:"tags"`
	} `json:"user_honour_info"`
	IsFollowed bool   `json:"is_followed"`
	TopPhoto   string `json:"top_photo"`
	Theme      struct {
	} `json:"theme"`
	SysNotice struct {
	} `json:"sys_notice"`
	LiveRoom struct {
		RoomStatus    int    `json:"roomStatus"`
		LiveStatus    int    `json:"liveStatus"`
		URL           string `json:"url"`
		Title         string `json:"title"`
		Cover         string `json:"cover"`
		Roomid        int    `json:"roomid"`
		RoundStatus   int    `json:"roundStatus"`
		BroadcastType int    `json:"broadcast_type"`
		WatchedShow   struct {
			Switch       bool   `json:"switch"`
			Num          int    `json:"num"`
			TextSmall    string `json:"text_small"`
			TextLarge    string `json:"text_large"`
			Icon         string `json:"icon"`
			IconLocation string `json:"icon_location"`
			IconWeb      string `json:"icon_web"`
		} `json:"watched_show"`
	} `json:"live_room"`
	Birthday string `json:"birthday"`
	School   struct {
		Name string `json:"name"`
	} `json:"school"`
	Profession struct {
		Name       string `json:"name"`
		Department string `json:"department"`
		Title      string `json:"title"`
		IsShow     int    `json:"is_show"`
	} `json:"profession"`
	Tags   []string `json:"tags"`
	Series struct {
		UserUpgradeStatus int  `json:"user_upgrade_status"`
		ShowUpgradeWindow bool `json:"show_upgrade_window"`
	} `json:"series"`
	IsSeniorMember int `json:"is_senior_member"`
}

func (a *API) GetUserInfo(uid int) (*UserInfoResp, error) {
	wts, wrid := a.GetWRID(map[string]interface{}{
		"mid": uid,
	})
	userInfo := &UserInfoResp{}
	_, err := a.CommonAPIClient.R().
		SetQueryParam("mid", strconv.Itoa(uid)).
		SetQueryParam("wts", strconv.FormatInt(wts, 10)).
		SetQueryParam("w_rid", wrid).
		SetResult(userInfo).
		Get("/x/space/wbi/acc/info")
	if err != nil {
		return nil, err
	}
	if userInfo.Code < 0 {
		return nil, errors.New("命中反爬策略")
	}
	return userInfo, nil
}

type FollowerInfoResp struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	TTL     int          `json:"ttl"`
	Data    FollowerInfo `json:"data"`
}

type FollowerInfo struct {
	Mid       int  `json:"mid"`
	Following uint `json:"following"`
	Follower  uint `json:"follower"`
}

func (a *API) GetFollowerInfo(uid int) (*FollowerInfoResp, error) {
	userInfo := &FollowerInfoResp{}
	_, err := a.CommonAPIClient.R().
		SetQueryParam("vmid", strconv.Itoa(uid)).
		SetResult(userInfo).
		Get("/x/relation/stat")
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}
