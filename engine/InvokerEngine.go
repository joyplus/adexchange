package engine

import (
	"admux/lib"
	m "admux/models"
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
var AdspaceSecretMap map[string]string

//key:<adspace_key>_<demand_id>; value:<demand_adspace_key>,<demand_secret_key>
var AdspaceMap map[string][]string

//key:<adspace_key>; value:<demand_id1>,<demand_id2>...
var AdspaceDemandMap map[string][]int

//key:<demand_id>; value:<demand_url>
var DemandMap map[int]string

//key:<adspace_key>_<demand_adspace_key>; value:<bool>
var AvbAdSpaceDemand map[string]bool

//key:<adspace_key>_<demand_adspace_key>; value:<bool>
var AvbAdspaceRegionTargeting map[string]bool

//key:<adspace_key>_<demand_adspace_key>_<region_code>; value:<left_imp>
var AvbAdSpaceRegion map[string]bool

func init() {
	test()
}

func test() {
	AdspaceMap = make(map[string][]string)
	AdspaceDemandMap = make(map[string][]int)
	DemandMap = make(map[int]string)

	AdspaceMap["TE57EAC5FA3FFACC_2"] = []string{"E757EAC5FA3FFACC", ""}
	AdspaceMap["TE57EAC5FA3FFACC_3"] = []string{"B4F1B7ABAA10D214", ""}
	AdspaceDemandMap["TE57EAC5FA3FFACC"] = []int{3, 2}
	DemandMap[2] = "http://ad.sandbox.madserving.com/adcall/bidrequest"
	DemandMap[3] = "http://api.sandbox.airwaveone.net/adcall/bidrequest"
}

//func generateParamsMapForMH(adRequest *m.AdRequest) map[string]string {
//	paramsMap := make(map[string]string)

//	return paramsMap
//}

func InvokeDemand(adRequest *m.AdRequest) *m.AdResponse {

	adspaceKey := adRequest.AdspaceId

	demandIds := AdspaceDemandMap[adspaceKey]

	if len(demandIds) == 0 {

		return &m.AdResponse{StatusCode: lib.ERROR_ILLEGAL_ADSPACE}
	}

	demandAry := make([]*Demand, len(demandIds))

	demandIndex := 0

	for _, demandId := range demandIds {
		demandUrl := DemandMap[demandId]

		key4AdspaceMap := adspaceKey + "_" + lib.ConvertIntToString(demandId)

		adspaceAry, ok := AdspaceMap[key4AdspaceMap]

		if ok {

			demand := new(Demand)
			demand.URL = demandUrl
			demand.AdRequest = adRequest
			demand.AdspaceKey = adspaceAry[0]
			demand.AdSecretKey = adspaceAry[1]
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
