package controllers

import (
	"adexchange/engine"
	"adexchange/lib"
	m "adexchange/models"
	"github.com/astaxie/beego"
	"time"
)

type WebviewRequestController struct {
	BaseController
}

//Request Ad for client
func (this *WebviewRequestController) WebviewReq() {
	t1 := time.Now().UnixNano()

	adRequest := m.AdRequest{}
	adResponse := new(m.AdResponse)
	beego.Debug(this.Ctx.Input.Request)
	if err := this.ParseForm(&adRequest); err != nil {

		adResponse.StatusCode = lib.ERROR_PARSE_REQUEST
	} else {
		adRequest.Did = lib.GenerateBid(adRequest.AdspaceKey)

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
			t2 := time.Now().UnixNano()
			adRequest.ProcessDuration = (t2 - t1) / 1000000
			engine.SendRequestLog(&adRequest, 1)
		}

	}

	beego.Debug(adRequest.AdspaceKey)
	tplName := engine.GetPmpAdspaceTemplate(adRequest.AdspaceKey)
	flg := engine.CheckTplName(tplName)

	beego.Debug(tplName)

	if flg {
		this.TplNames = "tpl/" + tplName + ".html"

		this.Data["statusCode"] = adResponse.StatusCode

		if adResponse.StatusCode == lib.STATUS_SUCCESS {
			adParam := map[string][]string{"clkTrackingUrls": adResponse.Adunit.ClkTrackingUrls, "implTrackingUrls": adResponse.Adunit.ImpTrackingUrls, "imgUrls": adResponse.Adunit.CreativeUrls}

			this.Data["AD"] = adParam
			this.Data["clickUrl"] = adResponse.Adunit.ClickUrl
			this.Data["width"] = adResponse.Adunit.AdWidth
			this.Data["height"] = adResponse.Adunit.AdHeight
		}

		this.Render()
	}

}
