package resource

type NavResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		WbiImg WbiImg `json:"wbi_img"`
	} `json:"data"`
}

type WbiImg struct {
	ImgUrl string `json:"img_url"`
	SubUrl string `json:"sub_url"`
}

func (a *API) Nav() (*NavResp, error) {
	navInfo := &NavResp{}
	_, err := a.CommonAPIClient.R().
		SetResult(navInfo).
		Get("/x/web-interface/nav")
	if err != nil {
		return nil, err
	}
	return navInfo, nil
}

func (a *API) GetWRID(params map[string]interface{}) (int64, string) {
	navInfo, _ := a.Nav()
	mainKey, subKey, _ := GetWbiKeys(navInfo.Data.WbiImg.ImgUrl, navInfo.Data.WbiImg.SubUrl)
	return EncWbi(params, mainKey, subKey)
}
