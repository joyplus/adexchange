package engine

import (
	"adexchange/lib"
	m "adexchange/models"
	//"bytes"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/franela/goreq"
	"net/url"
	"time"
)

func invokeMH(demand *Demand) {

	adRequest := demand.AdRequest
	beego.Debug("Start Invoke MH,bid:" + adRequest.Bid)
	item := url.Values{}

	//item.Set("bid", lib.GenerateBid(demand.AdspaceKey))
	item.Set("bid", adRequest.Bid)
	item.Set("adspaceid", demand.AdspaceKey)
	//hard code 2 to request MH as hero app
	item.Set("adtype", "2")
	item.Set("pkgname", adRequest.Pkgname)
	item.Set("appname", adRequest.Appname)
	item.Set("conn", adRequest.Conn)
	item.Set("carrier", adRequest.Carrier)
	//hard code 2 to return json response
	item.Set("apitype", "2")
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
	item.Set("lon", lib.ConvertFloatToString(adRequest.Lon))
	item.Set("lat", lib.ConvertFloatToString(adRequest.Lat))

	res, err := goreq.Request{
		Uri:         demand.URL,
		QueryString: item,
		//ShowDebug:   true,
		Timeout: time.Duration(demand.Timeout) * time.Millisecond,
	}.Do()

	adResponse := new(m.AdResponse)
	adResponse.Bid = adRequest.Bid
	adResponse.SetDemandAdspaceKey(demand.AdspaceKey)
	adResponse.SetResponseTime(time.Now().Unix())

	var strResponse string
	if serr, ok := err.(*goreq.Error); ok {
		beego.Critical(err.Error())
		if serr.Timeout() {
			adResponse = generateErrorResponse(adRequest, demand.AdspaceKey, lib.ERROR_TIMEOUT_ERROR)
			demand.Result <- adResponse
		} else {
			adResponse = generateErrorResponse(adRequest, demand.AdspaceKey, lib.ERROR_MHSERVER_ERROR)
			demand.Result <- adResponse
		}

	} else {
		var resultMap map[string]*m.MHAdUnit

		//flg, _ := beego.AppConfig.Bool("log_demand_body")
		//var err error
		//if flg {
		//	strResponse, _ = res.Body.ToString()
		//	err = json.Unmarshal([]byte(strResponse), &resultMap)

		//} else {
		//	err = res.Body.FromJsonTo(&resultMap)
		//}
		strResponse, _ = res.Body.ToString()
		err = json.Unmarshal([]byte(strResponse), &resultMap)
		defer res.Body.Close()

		if err != nil {
			beego.Critical(err.Error())
			adResponse = generateErrorResponse(adRequest, demand.AdspaceKey, lib.ERROR_MAP_ERROR)
			//demand.Result <- adResponse
		} else {
			if resultMap != nil {
				for _, v := range resultMap {
					adResponse = mapMHResult(v)
					adResponse.Bid = adRequest.Bid
					adResponse.SetDemandAdspaceKey(demand.AdspaceKey)
					//demand.Result <- adResponse
					break
				}
			} else {
				adResponse = generateErrorResponse(adRequest, demand.AdspaceKey, lib.ERROR_MAP_ERROR)
				//demand.Result <- adResponse
			}
		}
		if adResponse.StatusCode != lib.STATUS_SUCCESS {
			adResponse.ResBody = strResponse
			beego.Debug(adResponse.ResBody)
		}

		demand.Result <- adResponse

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
		//todo hardcode 3 for MH, only support picture ad
		//adUnit.CreativeType = 3
		adUnit.CreativeUrls = []string{mhAdunit.Imgurl}
		adUnit.ImpTrackingUrls = mhAdunit.Imgtracking
		adUnit.ClkTrackingUrls = mhAdunit.Thclkurl
		adUnit.AdWidth = mhAdunit.Adwidth
		adUnit.AdHeight = mhAdunit.Adheight
	}

	return adResponse
}
