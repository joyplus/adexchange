package controllers

import (
	"adexchange/engine"
	"adexchange/lib"
	m "adexchange/models"
	"github.com/astaxie/beego"
	"time"
)

type RequestController struct {
	BaseController
}

//Request Ad
func (this *RequestController) RequestAd() {
	adRequest := m.AdRequest{}
	adResponse := new(m.AdResponse)
	beego.Debug("Enter Request ad")
	if err := this.ParseForm(&adRequest); err != nil {

		adResponse.StatusCode = lib.ERROR_PARSE_REQUEST
	} else {
		adRequest.RequestTime = time.Now().Unix()
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
			SendLog(adRequest, 1)
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
	commonResponse := adResponse.GenerateCommonResponse()

	if adResponse.Adunit != nil && adResponse.Adunit.CreativeType == lib.CREATIVE_TYPE_HTML {
		cacheKey := lib.GetMd5String(adResponse.Bid)
		url := beego.AppConfig.String("viewad_server") + "?id=" + cacheKey
		commonResponse.SetHtmlCreativeUrl(url)
		SetCachedAdResponse(cacheKey, adResponse)
	}

	this.Data["json"] = commonResponse
	this.ServeJson()

}
