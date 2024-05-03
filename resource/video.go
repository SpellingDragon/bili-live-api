package resource

import "strconv"

// Pages
type Pages struct {
	Weblink   string    `json:"weblink"`
	Dimension Dimension `json:"dimension"`
	Cid       int       `json:"cid"`
	Page      int       `json:"page"`
	From      string    `json:"from"`
	Part      string    `json:"part"`
	Duration  int       `json:"duration"`
	Vid       string    `json:"vid"`
}

// Dimension
type Dimension struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	Rotate int `json:"rotate"`
}

// DescV2
type DescV2 struct {
	RawText string `json:"raw_text"`
	Type    int    `json:"type"`
	BizId   int    `json:"biz_id"`
}

// Owner
type Owner struct {
	Mid  int    `json:"mid"`
	Name string `json:"name"`
	Face string `json:"face"`
}

// UserGarb
type UserGarb struct {
	UrlImageAniCut string `json:"url_image_ani_cut"`
}

// DataDimension
type DataDimension struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	Rotate int `json:"rotate"`
}

// VideoInfoResponse
type VideoInfoResponse struct {
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    Video  `json:"data"`
	Code    int    `json:"code"`
}

// Video
type Video struct {
	Bvid               string        `json:"bvid"`
	Tname              string        `json:"tname"`
	Cid                int           `json:"cid"`
	Pages              []Pages       `json:"pages"`
	DescV2             []DescV2      `json:"desc_v2"`
	Duration           int           `json:"duration"`
	TeenageMode        int           `json:"teenage_mode"`
	NoCache            bool          `json:"no_cache"`
	Copyright          int           `json:"copyright"`
	Dynamic            string        `json:"dynamic"`
	IsStory            bool          `json:"is_story"`
	HonorReply         HonorReply    `json:"honor_reply"`
	Desc               string        `json:"desc"`
	State              int           `json:"state"`
	Owner              Owner         `json:"owner"`
	Videos             int           `json:"videos"`
	Rights             Rights        `json:"rights"`
	Tid                int           `json:"tid"`
	Pubdate            int           `json:"pubdate"`
	Ctime              int           `json:"ctime"`
	IsChargeableSeason bool          `json:"is_chargeable_season"`
	Subtitle           Subtitle      `json:"subtitle"`
	LikeIcon           string        `json:"like_icon"`
	Aid                int           `json:"aid"`
	Title              string        `json:"title"`
	Stat               Stat          `json:"stat"`
	Dimension          DataDimension `json:"dimension"`
	Premiere           interface{}   `json:"premiere"`
	UserGarb           UserGarb      `json:"user_garb"`
	Pic                string        `json:"pic"`
	IsSeasonDisplay    bool          `json:"is_season_display"`
}

// HonorReply
type HonorReply struct {
}

// Rights
type Rights struct {
	Elec          int `json:"elec"`
	Pay           int `json:"pay"`
	UgcPayPreview int `json:"ugc_pay_preview"`
	CleanMode     int `json:"clean_mode"`
	FreeWatch     int `json:"free_watch"`
	Bp            int `json:"bp"`
	Download      int `json:"download"`
	UgcPay        int `json:"ugc_pay"`
	IsSteinGate   int `json:"is_stein_gate"`
	Hd5           int `json:"hd5"`
	IsCooperation int `json:"is_cooperation"`
	NoBackground  int `json:"no_background"`
	ArcPay        int `json:"arc_pay"`
	Movie         int `json:"movie"`
	NoReprint     int `json:"no_reprint"`
	Autoplay      int `json:"autoplay"`
	Is360         int `json:"is_360"`
	NoShare       int `json:"no_share"`
}

// Subtitle
type Subtitle struct {
	AllowSubmit bool `json:"allow_submit"`
}

// Stat
type Stat struct {
	Coin       int    `json:"coin"`
	Share      int    `json:"share"`
	HisRank    int    `json:"his_rank"`
	Evaluation string `json:"evaluation"`
	Aid        int    `json:"aid"`
	View       int    `json:"view"`
	Favorite   int    `json:"favorite"`
	Like       int    `json:"like"`
	Dislike    int    `json:"dislike"`
	ArgueMsg   string `json:"argue_msg"`
	Danmaku    int    `json:"danmaku"`
	Reply      int    `json:"reply"`
	NowRank    int    `json:"now_rank"`
}

func (a *API) VideoInfo(bvId string) (*VideoInfoResponse, error) {
	videoInfo := &VideoInfoResponse{}
	_, err := a.CommonAPIClient.R().
		SetQueryParam("bvid", bvId).
		SetResult(videoInfo).
		Get("/x/web-interface/view")
	if err != nil {
		println(err)
		return nil, err
	}
	return videoInfo, nil
}

// VideoTagResponse
type VideoTagResponse struct {
	Code    int        `json:"code,omitempty"`
	Message string     `json:"message,omitempty"`
	Ttl     int        `json:"ttl,omitempty"`
	Data    []VideoTag `json:"data,omitempty"`
}

// VideoTag
type VideoTag struct {
	IsActivity      bool   `json:"is_activity"`
	Alpha           int    `json:"alpha"`
	TagId           int    `json:"tag_id"`
	Cover           string `json:"cover"`
	Liked           int    `json:"liked"`
	Ctime           int    `json:"ctime"`
	IsAtten         int    `json:"is_atten"`
	Likes           int    `json:"likes"`
	JumpUrl         string `json:"jump_url"`
	ShortContent    string `json:"short_content"`
	Attribute       int    `json:"attribute"`
	Color           string `json:"color"`
	SubscribedCount int    `json:"subscribed_count"`
	ArchiveCount    string `json:"archive_count"`
	State           int    `json:"state"`
	Hates           int    `json:"hates"`
	ExtraAttr       int    `json:"extra_attr"`
	MusicId         string `json:"music_id"`
	Content         string `json:"content"`
	Count           Count  `json:"count"`
	FeaturedCount   int    `json:"featured_count"`
	TagName         string `json:"tag_name"`
	IsSeason        bool   `json:"is_season"`
	HeadCover       string `json:"head_cover"`
	Hated           int    `json:"hated"`
	TagType         string `json:"tag_type"`
	Type            int    `json:"type"`
}

// Count
type Count struct {
	View  int `json:"view"`
	Use   int `json:"use"`
	Atten int `json:"atten"`
}

func (a *API) VideoTags(video Video) (*VideoTagResponse, error) {
	videoTags := &VideoTagResponse{}
	_, err := a.CommonAPIClient.R().
		SetQueryParam("aid", strconv.Itoa(video.Aid)).
		SetQueryParam("cid", strconv.Itoa(video.Cid)).
		SetResult(videoTags).
		Get("/x/web-interface/view/detail/tag")
	if err != nil {
		println(err)
		return nil, err
	}
	return videoTags, nil
}
