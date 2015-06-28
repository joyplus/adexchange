package models

//import (
//	"github.com/astaxie/beego"
//)

type AdResponse struct {
	StatusCode       int
	AdspaceKey       string
	DemandAdspaceKey string
	ResponseTime     int64
	Bid              string
	Adunit           *AdUnit
}

type CommonResponse struct {
	StatusCode int     `json:"statusCode"`
	AdspaceKey string  `json:"adspaceKey"`
	Bid        string  `json:"bid"`
	Adunit     *AdUnit `json:"adunit"`
}

func (this *AdResponse) GenerateCommonResponse() CommonResponse {
	res := CommonResponse{}
	res.StatusCode = this.StatusCode
	res.AdspaceKey = this.AdspaceKey
	res.Bid = this.Bid
	res.Adunit = this.Adunit

	return res
}

func (this *AdResponse) SetDemandAdspaceKey(dkey string) {
	this.DemandAdspaceKey = dkey
}

func (this *AdResponse) GetDemandAdspaceKey() string {
	return this.DemandAdspaceKey
}

func (this *AdResponse) SetResponseTime(responseTime int64) {
	this.ResponseTime = responseTime
}

func (this *AdResponse) GetResponseTime() int64 {
	return this.ResponseTime
}

func (this *AdResponse) AddImpTracking(url string) {
	if this.Adunit == nil {
		return
	}

	if this.Adunit.ImpTrackingUrls != nil {
		this.Adunit.ImpTrackingUrls = append(this.Adunit.ImpTrackingUrls, url)
	} else {
		this.Adunit.ImpTrackingUrls = []string{url}
	}

}

func (this *AdResponse) AddClkTracking(url string) {
	if this.Adunit == nil {
		return
	}

	if this.Adunit.ClkTrackingUrls != nil {
		this.Adunit.ClkTrackingUrls = append(this.Adunit.ClkTrackingUrls, url)
	} else {
		this.Adunit.ClkTrackingUrls = []string{url}
	}

}

type AdUnit struct {
	Cid             string   `json:"cid"`
	ClickUrl        string   `json:"clickUrl"`
	DisplayText     string   `json:"displayText"`
	ImageUrls       []string `json:"imageUrls"`
	ImpTrackingUrls []string `json:"impTrackingUrls"`
	ClkTrackingUrls []string `json:"clkTrackingUrls"`
	AdWidth         int      `json:"adWidth"`
	AdHeight        int      `json:"adHeight"`
}

type MHAdUnit struct {
	Adspaceid   string
	Returncode  int
	Cid         string
	Adwidth     int
	Adheight    int
	Adtype      int
	Imgurl      string
	Clickurl    string
	Imgtracking []string
	Thclkurl    []string
}

type AdspaceData struct {
	AdspaceKey string
	SecretKey  string
}

type PmpAdplaceInfo struct {
	PmpAdspaceKey    string
	DemandId         int
	DemandAdspaceKey string
	DemandSecretKey  string
}

type DemandInfo struct {
	DemandId           int
	Name               string
	RequestUrlTemplate string
	Timeout            int
}
