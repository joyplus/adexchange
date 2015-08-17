package engine

import (
	"adexchange/lib"
	m "adexchange/models"
	"bytes"
	"github.com/astaxie/beego"
	//"github.com/franela/goreq"
	//"net/url"
	"time"
)

//var c1 *httpclient.HttpClient

type Demand struct {
	URL           string
	Timeout       int
	AdRequest     *m.AdRequest
	AdspaceKey    string
	AdSecretKey   string
	Priority      int
	Result        chan *m.AdResponse
	TargetingCode string
	AppName       string
	PkgName       string
	Pcat          int
	Ua            string
}

//key:<adspace_key>; value:<PmpInfo>
var _PmpAdspaceMap map[string]m.PmpInfo

//key:<adspace_key>; value:<secret_key>
var _AdspaceSecretMap map[string]string

//key:<adspace_key>_<demand_id>; value:<demand_adspace_key>,<demand_secret_key>
var _AdspaceMap map[string]m.AdspaceData

//key:<adspace_key>; value:<demand_adspace_key1>,<demand_adspace_key2>...
var _AdspaceDemandMap map[string][]string

//key:<demand_id>; value:<demand_url>
var _DemandMap map[int]m.DemandInfo

//key:<adspace_key>_<demand_adspace_key>; value:<*AvbDemand>
var _AvbAdspaceDemand map[string]*m.AvbDemand

var _FuncMap lib.Funcs

var IMP_TRACKING_SERVER string
var CLK_TRACKING_SERVER string

func init() {
	_FuncMap = lib.NewFuncs(1)
	err := _FuncMap.Bind("invokeMH", invokeMH)
	if err != nil {
		beego.Emergency(err.Error())
	}
	err = _FuncMap.Bind("invokeCampaign", invokeCampaign)
	if err != nil {
		beego.Emergency(err.Error())
	}
	err = _FuncMap.Bind("invokeBD", invokeBD)
	if err != nil {
		beego.Emergency(err.Error())
	}
	err = _FuncMap.Bind("invokeMHQueue", invokeMHQueue)
	if err != nil {
		beego.Emergency(err.Error())
	}

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

	if _AdspaceMap == nil || _AdspaceDemandMap == nil || _DemandMap == nil || _PmpAdspaceMap == nil {
		return generateErrorResponse(adRequest, "", lib.ERROR_INITIAL_FAILED)
	}

	adspaceKey := adRequest.AdspaceKey
	if _, ok := _PmpAdspaceMap[adspaceKey]; !ok {
		return generateErrorResponse(adRequest, "", lib.ERROR_NO_PMP_ADSPACE_ERROR)
	}

	aryDemandAdspaceKeys := _AdspaceDemandMap[adspaceKey]
	beego.Debug(aryDemandAdspaceKeys)
	if len(aryDemandAdspaceKeys) == 0 {

		return generateErrorResponse(adRequest, "", lib.ERROR_ILLEGAL_ADSPACE)
	}

	demandAry := make([]*Demand, len(aryDemandAdspaceKeys))

	demandIndex := 0

	for _, demandAdspaceKey := range aryDemandAdspaceKeys {

		key4AdspaceMap := adspaceKey + "_" + demandAdspaceKey
		//beego.Debug(key4AdspaceMap)
		adspaceData, ok := _AdspaceMap[key4AdspaceMap]

		avbFlg, targetingCode := checkAvbDemand(adRequest, adspaceData)

		if ok && avbFlg {

			demandInfo := _DemandMap[adspaceData.DemandId]

			demand := new(Demand)
			demand.URL = demandInfo.RequestUrlTemplate
			demand.Timeout = demandInfo.Timeout
			demand.AdRequest = adRequest
			demand.AdspaceKey = adspaceData.AdspaceKey
			demand.AdSecretKey = adspaceData.SecretKey
			demand.Priority = adspaceData.Priority
			demand.TargetingCode = targetingCode
			demand.Result = make(chan *m.AdResponse)

			//mockup app info
			beego.Debug(adspaceData)
			demand.AppName = adspaceData.AppName
			demand.PkgName = adspaceData.PkgName
			demand.Pcat = adspaceData.Pcat
			demand.Ua = adspaceData.Ua
			demandAry[demandIndex] = demand
			demandIndex++
			//go invokeMH(demand)
			go _FuncMap.Call(demandInfo.InvokeFuncName, demand)
		}
	}

	if demandIndex == 0 {
		return nil
	}

	adResultAry := make([]*m.AdResponse, demandIndex)
	successIndex := 0
	var tmp *m.AdResponse
	for index := 0; index < demandIndex; index++ {
		demand := demandAry[index]
		tmp = <-demand.Result
		if tmp != nil && tmp.StatusCode == 200 {
			adResultAry[successIndex] = tmp
			successIndex++
		}

		go SendDemandLog(tmp)
	}

	if successIndex == 0 {
		return tmp
	}

	beego.Debug(successIndex)

	adResponse := chooseAdResponse(adResultAry[:successIndex])
	adResponse.AdspaceKey = adRequest.AdspaceKey
	adRequest.DemandAdspaceKey = adResponse.DemandAdspaceKey
	if adResponse.StatusCode == 200 {
		adResponse.Adunit.CreativeType = _PmpAdspaceMap[adResponse.AdspaceKey].CreativeType
		impTrackUrl, clkTrackUrl := generateTrackingUrl(adResponse, adRequest)
		adResponse.AddImpTracking(impTrackUrl)
		adResponse.PmpClkTrackingUrl = clkTrackUrl
	}

	return adResponse
}

func generateTrackingUrl(adResponse *m.AdResponse, adRequest *m.AdRequest) (string, string) {
	var buffer bytes.Buffer
	buffer.WriteString("bid=")
	buffer.WriteString(adResponse.Bid)
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

func chooseAdResponse(aryAdResponse []*m.AdResponse) (adResponse *m.AdResponse) {

	//for _, adResponse = range aryAdResponse {
	//	if adResponse != nil && adResponse.StatusCode == 200 {
	//		return adResponse
	//		break
	//	}
	//}
	if len(aryAdResponse) == 1 {

		adResponse = aryAdResponse[0]

	} else if len(aryAdResponse) > 1 {

		aryWeightItem := make([]*lib.WeightItem, len(aryAdResponse))
		currentIndex := 0
		selectedIndex := 999

		for i, adResponseItem := range aryAdResponse {
			if adResponseItem.Priority >= 100 {
				selectedIndex = i
				break
			}
			weightItem := new(lib.WeightItem)
			aryWeightItem[i] = weightItem
			weightItem.Weight = adResponseItem.Priority
			beego.Debug(adResponseItem.Priority)
			weightItem.StartNumber = currentIndex + 1
			weightItem.EndNumber = currentIndex + weightItem.Weight
			weightItem.Index = i
			currentIndex = currentIndex + weightItem.Weight

		}

		if selectedIndex == 999 {
			selectedIndex = lib.ChooseItem(aryWeightItem)
		}

		beego.Debug("SelectedIndex: " + lib.ConvertIntToString(selectedIndex))

		adResponse = aryAdResponse[selectedIndex]

		beego.Debug(aryAdResponse)
		//random := lib.GetRandomNumber(0, len(aryAdResponse))
		//adResponse = aryAdResponse[random]
	}

	return
}

func SetupAdspaceMap(adspaceMap map[string]m.AdspaceData) {
	_AdspaceMap = adspaceMap
}
func SetupAdspaceDemandMap(adspaceDemandMap map[string][]string) {
	_AdspaceDemandMap = adspaceDemandMap
}
func SetupDemandMap(demandMap map[int]m.DemandInfo) {
	_DemandMap = demandMap
}
func SetupAvbAdspaceDemandMap(avbDemandMap map[string]*m.AvbDemand) {
	_AvbAdspaceDemand = avbDemandMap
}
func SetupPmpAdspaceMap(pmpAdspaceMap map[string]m.PmpInfo) {
	_PmpAdspaceMap = pmpAdspaceMap
}

func checkAvbDemand(adRequest *m.AdRequest, adspaceData m.AdspaceData) (avbFlg bool, targetingCode string) {

	beego.Debug("Start to Check avb demand")
	key := adRequest.AdspaceKey + "_" + adspaceData.AdspaceKey

	beego.Debug("avb key:" + key)
	if avbDemand, ok := _AvbAdspaceDemand[key]; ok {
		avbFlg, targetingCode = avbDemand.CheckAvailable(adRequest)
	}

	return avbFlg, targetingCode

}

func generateErrorResponse(adRequest *m.AdRequest, demandAdspaceKey string, statusCode int) *m.AdResponse {
	adResponse := new(m.AdResponse)
	adResponse.StatusCode = statusCode
	adResponse.Bid = adRequest.Bid
	adResponse.AdspaceKey = adRequest.AdspaceKey
	adResponse.DemandAdspaceKey = demandAdspaceKey
	adResponse.ResponseTime = time.Now().Unix()
	return adResponse
}

func initAdResponse(demand *Demand) (adResponse *m.AdResponse) {
	adResponse = new(m.AdResponse)
	adResponse.Bid = demand.AdRequest.Bid
	adResponse.AdspaceKey = demand.AdRequest.AdspaceKey
	adResponse.SetDemandAdspaceKey(demand.AdspaceKey)
	adResponse.SetResponseTime(time.Now().Unix())
	adResponse.Priority = demand.Priority

	return adResponse
}
