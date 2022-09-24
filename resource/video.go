package resource

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
	Data    Data   `json:"data"`
	Code    int    `json:"code"`
}

// Data
type Data struct {
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

func VideoInfo(bvId string) (*VideoInfoResponse, error) {
	videoInfo := &VideoInfoResponse{}
	_, err := apiClient.R().
		SetQueryParam("bvid", bvId).
		SetResult(videoInfo).
		Get("/x/web-interface")
	if err != nil {
		return nil, err
	}
	return videoInfo, nil
}
