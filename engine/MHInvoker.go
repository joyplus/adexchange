package engine

import (
	"adexchange/lib"
	m "adexchange/models"
	//"bytes"
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
