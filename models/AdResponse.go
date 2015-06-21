package models

type AdResponse struct {
	StatusCode int
	Bid        string
	Adunit     *AdUnit
}

type AdUnit struct {
	Cid             string
	ClickUrl        string
	DisplayText     string
	ImageUrls       []string
	ImpTrackingUrls []string
	ClkTrackingUrls []string
	AdWidth         int
	AdHeight        int
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
	DemandId int
	Url      string
}
