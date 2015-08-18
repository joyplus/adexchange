package controllers

import (
	"adexchange/engine"
	"adexchange/lib"
	m "adexchange/models"
	"github.com/astaxie/beego"
	"time"
)

type ClientRequestController struct {
	BaseController
}

//Request Ad for client
func (this *ClientRequestController) RequestAd4Client() {
	adRequest := m.AdRequest{}
	adResponse := new(m.AdResponse)
	beego.Debug(this.Ctx.Input.Request)
	if err := this.ParseForm(&adRequest); err != nil {

		adResponse.StatusCode = lib.ERROR_PARSE_REQUEST
	} else {
		adRequest.Did = GenerateBid(adRequest)

		adRequest.RequestTime = time.Now().Unix()
		clientIp := GetClientIP(this.Ctx.Input)
		beego.Debug("Request Client IP:" + clientIp)
		adRequest.Ip = clientIp

		ua := this.Ctx.Input.Header("User-Agent")
		beego.Debug("Request UA:" + ua)
		adRequest.Ua = ua

		tmp := engine.InvokeDemand(&adRequest)

		if tmp == nil {
			adResponse.StatusCode = lib.ERROR_NO_DEMAND_ERROR
			adResponse.Bid = adRequest.Bid
			adResponse.AdspaceKey = adRequest.AdspaceKey
		} else {
			adResponse = tmp
		}

		//only running pmp adspace need track request log
		if adResponse.StatusCode != lib.ERROR_NO_PMP_ADSPACE_ERROR {
			adRequest.StatusCode = adResponse.StatusCode
			go SendLog(adRequest, 1)
		}

		//if err != nil {
		//	beego.Debug("Enter sss ad")
		//	if e, ok := err.(*lib.SysError); ok {
		//		adResponse.StatusCode = e.ErrorCode
		//	} else {
		//		adResponse.StatusCode = lib.ERROR_UNKNON_ERROR
		//	}
		//	beego.Debug("Enter ssaass ad")
		//}
	}
	//commonResponse := adResponse.GenerateCommonResponse()

	//if adResponse.Adunit != nil {
	//	if adResponse.Adunit.CreativeType == lib.CREATIVE_TYPE_HTML {
	//		cacheKey := lib.GetMd5String(adResponse.Bid)
	//		url := beego.AppConfig.String("viewad_server") + "?id=" + cacheKey
	//		commonResponse.SetHtmlCreativeUrl(url)
	//		SetCachedAdResponse(cacheKey, adResponse)
	//	} else {
	//		cacheKey := lib.GetMd5String(adResponse.Bid)
	//		SetCachedClkUrl(cacheKey, adResponse.Adunit.ClickUrl)
	//		adResponse.Adunit.ClickUrl = adResponse.PmpClkTrackingUrl
	//	}
	//}

	commonResponse := GetCommonResponse(adResponse)

	this.Data["json"] = commonResponse
	this.ServeJson()

}
