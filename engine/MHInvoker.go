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
	beego.Debug("Start Invoke MH,did:" + adRequest.Did)
	item := url.Values{}

	//item.Set("bid", lib.GenerateBid(demand.AdspaceKey))
	item.Set("bid", demand.Did)
	item.Set("adspaceid", demand.AdspaceKey)
	//hard code 2 to request MH as hero app
	item.Set("adtype", "2")

	if len(demand.PkgName) > 0 {
		item.Set("pkgname", demand.PkgName)
	} else {
		item.Set("pkgname", adRequest.Pkgname)
	}

	if len(demand.PkgName) > 0 {
		item.Set("appname", demand.AppName)
	} else {
		item.Set("appname", adRequest.Appname)
	}

	if len(demand.PkgName) > 0 {
		item.Set("pcat", lib.ConvertIntToString(demand.Pcat))
	} else {
		item.Set("pcat", adRequest.Pcat)
	}

	if len(demand.PkgName) > 0 {
		item.Set("ua", demand.Ua)
	} else {
		item.Set("ua", adRequest.Ua)
	}

	item.Set("conn", adRequest.Conn)
	item.Set("carrier", adRequest.Carrier)
	//hard code 2 to return json response
	//hard code 4 to return json response with display titile and text
	item.Set("apitype", "4")
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

	adResponse := initAdResponse(demand)

	var strResponse string
	if serr, ok := err.(*goreq.Error); ok {
		beego.Critical(err.Error())
		if serr.Timeout() {
			//adResponse = generateErrorResponse(adRequest, demand.AdspaceKey, lib.ERROR_TIMEOUT_ERROR)
			adResponse.StatusCode = lib.ERROR_TIMEOUT_ERROR
		} else {
			//adResponse = generateErrorResponse(adRequest, demand.AdspaceKey, lib.ERROR_MHSERVER_ERROR)
			adResponse.StatusCode = lib.ERROR_MHSERVER_ERROR
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
		err = json.Unmarshal(lib.EscapeCtrl([]byte(strResponse)), &resultMap)
		defer res.Body.Close()

		if err != nil {
			beego.Critical(err.Error())
			//adResponse = generateErrorResponse(adRequest, demand.AdspaceKey, lib.ERROR_MAP_ERROR)
			adResponse.StatusCode = lib.ERROR_MAP_ERROR
			//demand.Result <- adResponse
		} else {
			if resultMap != nil {
				for _, v := range resultMap {
					mapMHResult(adResponse, v)
					//adResponse.Bid = adRequest.Bid
					//adResponse.SetDemandAdspaceKey(demand.AdspaceKey)
					//demand.Result <- adResponse
					break
				}
			} else {
				//adResponse = generateErrorResponse(adRequest, demand.AdspaceKey, lib.ERROR_MAP_ERROR)
				//demand.Result <- adResponse
				adResponse.StatusCode = lib.ERROR_MAP_ERROR
			}
		}
	}

	if adResponse.StatusCode != lib.STATUS_SUCCESS {
		adResponse.ResBody = strResponse
	}

	go SendDemandLog(adResponse)

	demand.Result <- adResponse
}

func mapMHResult(adResponse *m.AdResponse, mhAdunit *m.MHAdUnit) {

	adResponse.StatusCode = mhAdunit.Returncode
	adResponse.Did = mhAdunit.Bid

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
		adUnit.DisplayTitle = mhAdunit.Displaytitle
		adUnit.DisplayText = mhAdunit.Displaytext
	}

	//return adResponse
}
