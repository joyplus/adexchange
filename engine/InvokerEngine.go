package engine

import (
	"adexchange/lib"
	m "adexchange/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"time"
)

//var c1 *httpclient.HttpClient

type Demand struct {
	URL         string
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
var _DemandMap map[int]string

//key:<adspace_key>_<demand_adspace_key>; value:<bool>
var _AvbAdSpaceDemand map[string]bool

//key:<adspace_key>_<demand_adspace_key>; value:<bool>
var _AvbAdspaceRegionTargeting map[string]bool

//key:<adspace_key>_<demand_adspace_key>_<region_code>; value:<left_imp>
var _AvbAdSpaceRegion map[string]bool

func init() {
	test()
}

func test() {
	_AdspaceMap = make(map[string]m.AdspaceData)
	_AdspaceDemandMap = make(map[string][]int)
	_DemandMap = make(map[int]string)

	_AdspaceMap["TE57EAC5FA3FFACC_2"] = m.AdspaceData{AdspaceKey: "E757EAC5FA3FFACC"}
	_AdspaceMap["TE57EAC5FA3FFACC_3"] = m.AdspaceData{AdspaceKey: "B4F1B7ABAA10D214"}
	_AdspaceDemandMap["TE57EAC5FA3FFACC"] = []int{3, 2}
	_DemandMap[2] = "http://ad.sandbox.madserving.com/adcall/bidrequest"
	_DemandMap[3] = "http://api.sandbox.airwaveone.net/adcall/bidrequest"
}

//func generateParamsMapForMH(adRequest *m.AdRequest) map[string]string {
//	paramsMap := make(map[string]string)

//	return paramsMap
//}

func InvokeDemand(adRequest *m.AdRequest) *m.AdResponse {

	adspaceKey := adRequest.AdspaceKey

	demandIds := _AdspaceDemandMap[adspaceKey]

	if len(demandIds) == 0 {

		return &m.AdResponse{StatusCode: lib.ERROR_ILLEGAL_ADSPACE}
	}

	demandAry := make([]*Demand, len(demandIds))

	demandIndex := 0

	for _, demandId := range demandIds {
		demandUrl := _DemandMap[demandId]

		key4AdspaceMap := adspaceKey + "_" + lib.ConvertIntToString(demandId)

		adspaceData, ok := _AdspaceMap[key4AdspaceMap]

		if ok {

			demand := new(Demand)
			demand.URL = demandUrl
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
	}

	//for index, demand := range demandAry {
	//	if demand != nil {
	//		tmp := <-demand.Result
	//		adResultAry[index] = tmp
	//	}

	//}

	return chooseAdResponse(adResultAry)
}

func invokeMH(demand *Demand) {

	beego.Debug("Start Invoke MH")
	req := httplib.Get(demand.URL).Debug(true).SetTimeout(400*time.Millisecond, 300*time.Millisecond)

	adRequest := demand.AdRequest
	req.Param("bid", lib.GenerateBid(demand.AdspaceKey))
	req.Param("adspaceid", demand.AdspaceKey)
	req.Param("adtype", adRequest.AdType)
	req.Param("pkgname", adRequest.Pkgname)
	req.Param("appname", adRequest.Appname)
	req.Param("conn", adRequest.Conn)
	req.Param("carrier", adRequest.Carrier)
	req.Param("apitype", adRequest.ApiType)
	req.Param("os", adRequest.Os)
	req.Param("osv", adRequest.Osv)
	req.Param("imei", adRequest.Imei)
	req.Param("wma", adRequest.Wma)
	req.Param("aid", adRequest.Aid)
	req.Param("aaid", adRequest.Aaid)
	req.Param("idfa", adRequest.Idfa)
	req.Param("oid", adRequest.Oid)
	req.Param("uid", adRequest.Uid)
	req.Param("device", adRequest.Device)
	req.Param("ua", adRequest.Ua)
	req.Param("ip", adRequest.Ip)
	req.Param("width", adRequest.Width)
	req.Param("height", adRequest.Height)
	req.Param("density", adRequest.Ua)
	req.Param("lon", adRequest.Lon)
	req.Param("lat", adRequest.Lat)

	var resultMap map[string]*m.MHAdUnit
	req.ToJson(&resultMap)

	if resultMap != nil {
		for _, v := range resultMap {
			demand.Result <- mapMHResult(v)
			break
		}
	} else {
		demand.Result <- generateErrorResponse(lib.ERROR_MH_ERROR)
	}

}

func mapMHResult(mhAdunit *m.MHAdUnit) (adResponse *m.AdResponse) {

	adResponse = new(m.AdResponse)
	adResponse.StatusCode = mhAdunit.Returncode

	beego.Debug(mhAdunit)
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
func SetupDemandMap(demandMap map[int]string) {
	_DemandMap = demandMap
}
