package models

//import (
//	"github.com/astaxie/beego"
//)

type AdResponse struct {
	StatusCode       int    `json:"statusCode"`
	AdspaceKey       string `json:"adspaceKey"`
	demandAdspaceKey string
	Bid              string  `json:"bid"`
	Adunit           *AdUnit `json:"adunit"`
}

func (this *AdResponse) SetDemandAdspaceKey(dkey string) {
	this.demandAdspaceKey = dkey
}

func (this *AdResponse) GetDemandAdspaceKey() string {
	return this.demandAdspaceKey
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
