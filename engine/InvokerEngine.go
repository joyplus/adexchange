package engine

import (
	"adexchange/lib"
	m "adexchange/models"
	"bytes"
	"github.com/astaxie/beego"
	"github.com/franela/goreq"
	"net/url"
	"time"
)

//var c1 *httpclient.HttpClient

type Demand struct {
	URL         string
	Timeout     int
	AdRequest   *m.AdRequest
	AdspaceKey  string
	AdSecretKey string
	Result      chan *m.AdResponse
}

//key:<adspace_key>; value:<secret_key>
var _AdspaceSecretMap map[string]string

//key:<adspace_key>_<demand_id>; value:<demand_adspace_key>,<demand_secret_key>
var _AdspaceMap map[string]m.AdspaceData

//key:<adspace_key>; value:<demand_id1>,<demand_id2>...
var _AdspaceDemandMap map[string][]int

//key:<demand_id>; value:<demand_url>
var _DemandMap map[int]m.DemandInfo

//key:<adspace_key>_<demand_adspace_key>; value:<bool>
var _AvbAdSpaceDemand map[string]bool

//key:<adspace_key>_<demand_adspace_key>; value:<bool>
var _AvbAdspaceRegionTargeting map[string]bool

//key:<adspace_key>_<demand_adspace_key>_<region_code>; value:<left_imp>
var _AvbAdSpaceRegion map[string]bool

var IMP_TRACKING_SERVER string
var CLK_TRACKING_SERVER string

func init() {

	IMP_TRACKING_SERVER = beego.AppConfig.String("imp_tracking_server")
	CLK_TRACKING_SERVER = beego.AppConfig.String("clk_tracking_server")

}

//func test() {
//	_AdspaceMap = make(map[string]m.AdspaceData)
//	_AdspaceDemandMap = make(map[string][]int)
//	_DemandMap = make(map[int]string)

//	_AdspaceMap["TE57EAC5FA3FFACC_2"] = m.AdspaceData{AdspaceKey: "E757EAC5FA3FFACC"}
//	_AdspaceMap["TE57EAC5FA3FFACC_3"] = m.AdspaceData{AdspaceKey: "B4F1B7ABAA10D214"}
//	_AdspaceDemandMap["TE57EAC5FA3FFACC"] = []int{3, 2}
//	_DemandMap[2] = "http://ad.sandbox.madserving.com/adcall/bidrequest"
//	_DemandMap[3] = "http://api.sandbox.airwaveone.net/adcall/bidrequest"
//}

//func generateParamsMapForMH(adRequest *m.AdRequest) map[string]string {
//	paramsMap := make(map[string]string)

//	return paramsMap
//}

func InvokeDemand(adRequest *m.AdRequest) *m.AdResponse {

	if _AdspaceMap == nil || _AdspaceDemandMap == nil || _DemandMap == nil {
		return &m.AdResponse{StatusCode: lib.ERROR_INITIAL_FAILED}
	}

	adspaceKey := adRequest.AdspaceKey

	demandIds := _AdspaceDemandMap[adspaceKey]

	if len(demandIds) == 0 {

		return &m.AdResponse{StatusCode: lib.ERROR_ILLEGAL_ADSPACE}
	}

	demandAry := make([]*Demand, len(demandIds))

	demandIndex := 0

	for _, demandId := range demandIds {

		key4AdspaceMap := adspaceKey + "_" + lib.ConvertIntToString(demandId)

		adspaceData, ok := _AdspaceMap[key4AdspaceMap]

		if ok {

			demandInfo := _DemandMap[demandId]

			demand := new(Demand)
			demand.URL = demandInfo.RequestUrlTemplate
			demand.Timeout = demandInfo.Timeout
			demand.AdRequest = adRequest
			demand.AdspaceKey = adspaceData.AdspaceKey
			demand.AdSecretKey = adspaceData.SecretKey
			demand.Result = make(chan *m.AdResponse)
			demandAry[demandIndex] = demand
			demandIndex++
			go invokeMH(demand)
		}
	}

	adResultAry := make([]*m.AdResponse, demandIndex)

	for index := 0; index < demandIndex; index++ {
		demand := demandAry[index]
		tmp := <-demand.Result
		adResultAry[index] = tmp

		SendDemandLog(tmp)
	}

	//for index, demand := range demandAry {
	//	if demand != nil {
	//		tmp := <-demand.Result
	//		adResultAry[index] = tmp
	//	}

	//}
	adResponse := chooseAdResponse(adResultAry)
	adResponse.AdspaceKey = adRequest.AdspaceKey
	if adResponse.StatusCode == 200 {
		impTrackUrl, clkTrackUrl := generateTrackingUrl(adRequest)
		adResponse.AddImpTracking(impTrackUrl)
		adResponse.AddClkTracking(clkTrackUrl)
	}

	return adResponse
}

func generateTrackingUrl(adRequest *m.AdRequest) (string, string) {
	var buffer bytes.Buffer
	buffer.WriteString("bid=")
	buffer.WriteString(adRequest.Bid)
	buffer.WriteString("&adspaceid=")
	buffer.WriteString(adRequest.AdspaceKey)
	buffer.WriteString("&dkey=")
	buffer.WriteString(adRequest.DemandAdspaceKey)
	buffer.WriteString("&pkgname=")
	buffer.WriteString(adRequest.Pkgname)
	buffer.WriteString("&os=")
	buffer.WriteString(lib.ConvertIntToString(adRequest.Os))
	buffer.WriteString("&imei=")
	buffer.WriteString(adRequest.Imei)
	buffer.WriteString("&wma=")
	buffer.WriteString(adRequest.Wma)
	buffer.WriteString("&aid=")
	buffer.WriteString(adRequest.Aid)
	buffer.WriteString("&aaid=")
	buffer.WriteString(adRequest.Aaid)
	buffer.WriteString("&idfa=")
	buffer.WriteString(adRequest.Idfa)
	buffer.WriteString("&oid=")
	buffer.WriteString(adRequest.Oid)
	buffer.WriteString("&uid=")
	buffer.WriteString(adRequest.Uid)
	buffer.WriteString("&ua=")
	buffer.WriteString(adRequest.Ua)

	paramStr := buffer.String()
	impTrackUrl := IMP_TRACKING_SERVER + "?" + paramStr
	clkTrackUrl := CLK_TRACKING_SERVER + "?" + paramStr

	return impTrackUrl, clkTrackUrl

}

func invokeMH(demand *Demand) {

	beego.Debug("Start Invoke MH")
	//	req := httplib.Get(demand.URL).Debug(true).SetTimeout(400*time.Millisecond, 300*time.Millisecond)

	adRequest := demand.AdRequest
	item := url.Values{}

	//item.Set("bid", lib.GenerateBid(demand.AdspaceKey))
	item.Set("bid", adRequest.Bid)
	item.Set("adspaceid", demand.AdspaceKey)
	item.Set("adtype", adRequest.AdType)
	item.Set("pkgname", adRequest.Pkgname)
	item.Set("appname", adRequest.Appname)
	item.Set("conn", adRequest.Conn)
	item.Set("carrier", adRequest.Carrier)
	item.Set("apitype", adRequest.ApiType)
	item.Set("os", lib.ConvertIntToString(adRequest.Os))
	item.Set("osv", adRequest.Osv)
	item.Set("imei", adRequest.Imei)
	item.Set("wma", adRequest.Wma)
	item.Set("aid", adRequest.Aid)
	item.Set("aaid", adRequest.Aaid)
	item.Set("idfa", adRequest.Idfa)
	item.Set("oid", adRequest.Oid)
	item.Set("uid", adRequest.Uid)
	item.Set("device", adRequest.Device)
	item.Set("ua", adRequest.Ua)
	item.Set("ip", adRequest.Ip)
	item.Set("width", adRequest.Width)
	item.Set("height", adRequest.Height)
	item.Set("density", adRequest.Density)
	item.Set("lon", adRequest.Lon)
	item.Set("lat", adRequest.Lat)

	res, err := goreq.Request{
		Uri:         demand.URL,
		QueryString: item,
		Timeout:     time.Duration(demand.Timeout) * time.Millisecond,
	}.Do()

	adResponse := new(m.AdResponse)
	adResponse.Bid = adRequest.Bid
	adResponse.SetDemandAdspaceKey(demand.AdspaceKey)

	if serr, ok := err.(*goreq.Error); ok {
		beego.Error(err.Error())
		if serr.Timeout() {
			adResponse.StatusCode = lib.ERROR_TIMEOUT_ERROR
			demand.Result <- adResponse
		} else {
			adResponse.StatusCode = lib.ERROR_MHSERVER_ERROR
			demand.Result <- adResponse
		}
	} else {
		var resultMap map[string]*m.MHAdUnit

		err = res.Body.FromJsonTo(&resultMap)

		defer res.Body.Close()

		if err != nil {
			beego.Error(err.Error())
			adResponse.StatusCode = lib.ERROR_MAP_ERROR
			demand.Result <- adResponse
		} else {
			if resultMap != nil {
				for _, v := range resultMap {
					adResponse = mapMHResult(v)
					adResponse.Bid = adRequest.Bid
					adResponse.SetDemandAdspaceKey(demand.AdspaceKey)
					demand.Result <- adResponse
					break
				}
			} else {
				adResponse.StatusCode = lib.ERROR_MAP_ERROR
				demand.Result <- adResponse
			}
		}

	}

}

func mapMHResult(mhAdunit *m.MHAdUnit) (adResponse *m.AdResponse) {

	adResponse = new(m.AdResponse)
	adResponse.StatusCode = mhAdunit.Returncode
	adResponse.SetResponseTime(time.Now().Unix())

	if adResponse.StatusCode == 200 {
		adUnit := new(m.AdUnit)
		adResponse.Adunit = adUnit
		adUnit.Cid = mhAdunit.Cid
		adUnit.ClickUrl = mhAdunit.Clickurl
		adUnit.ImageUrls = []string{mhAdunit.Imgurl}
		adUnit.ImpTrackingUrls = mhAdunit.Imgtracking
		adUnit.ClkTrackingUrls = mhAdunit.Thclkurl
		adUnit.AdWidth = mhAdunit.Adwidth
		adUnit.AdHeight = mhAdunit.Adheight
	}

	return adResponse
}

func chooseAdResponse(aryAdResponse []*m.AdResponse) (adResponse *m.AdResponse) {

	for _, adResponse = range aryAdResponse {
		if adResponse != nil && adResponse.StatusCode == 200 {
			return adResponse
			break
		}
	}

	return adResponse
}

func generateErrorResponse(statusCode int) (adResponse *m.AdResponse) {
	adResponse = new(m.AdResponse)
	adResponse.StatusCode = statusCode

	return adResponse
}

func UpdateAdspaceStatus(adspaceKey string, demandAdspaceKey string, status bool) {
	_AvbAdSpaceDemand[adspaceKey+"_"+demandAdspaceKey] = status

}

func SetupAdspaceSecretMap(adspaceSecretMap map[string]string) {
	_AdspaceSecretMap = adspaceSecretMap
}
func SetupAdspaceMap(adspaceMap map[string]m.AdspaceData) {
	_AdspaceMap = adspaceMap
}
func SetupAdspaceDemandMap(adspaceDemandMap map[string][]int) {
	_AdspaceDemandMap = adspaceDemandMap
}
func SetupDemandMap(demandMap map[int]m.DemandInfo) {
	_DemandMap = demandMap
}
