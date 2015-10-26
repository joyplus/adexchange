package controllers

import (
	"adexchange/engine"
	"adexchange/lib"
	m "adexchange/models"
	"github.com/astaxie/beego"
	"time"
)

type TrackingController struct {
	beego.Controller
}

func (this *RequestController) TrackImp() {
	adRequest := m.AdRequest{}
	//adResponse := new(m.TrackResponse)
	beego.Debug("Enter Tracking imp")

	if err := this.ParseForm(&adRequest); err != nil {
		adRequest.StatusCode = lib.ERROR_PARSE_REQUEST
	} else {
		adRequest.StatusCode = lib.STATUS_SUCCESS
		clientIp := GetClientIP(this.Ctx.Input)
		beego.Debug("Clk Client IP:" + clientIp)
		adRequest.Ip = clientIp
		adRequest.RequestTime = time.Now().Unix()
		ua := this.Ctx.Input.Header("User-Agent")
		adRequest.Ua = ua
		engine.SendRequestLog(&adRequest, 2)
	}

	this.Ctx.Redirect(302, beego.AppConfig.String("public_server")+"/1.gif")
	//this.Data["json"] = &adResponse
	//this.ServeJson()

}

func (this *RequestController) TrackClk() {
	adRequest := m.AdRequest{}
	//adResponse := new(m.TrackResponse)
	beego.Debug("Enter Tracking clk")

	if err := this.ParseForm(&adRequest); err != nil {
		adRequest.StatusCode = lib.ERROR_PARSE_REQUEST
	} else {
		adRequest.StatusCode = lib.STATUS_SUCCESS
		clientIp := GetClientIP(this.Ctx.Input)
		beego.Debug("Imp Client IP:" + clientIp)
		adRequest.Ip = clientIp
		adRequest.RequestTime = time.Now().Unix()
		ua := this.Ctx.Input.Header("User-Agent")
		adRequest.Ua = ua
		engine.SendRequestLog(&adRequest, 3)
	}

	cacheKey := adRequest.Cuk

	originalTrackingUrl := GetCachedClkUrl(cacheKey)

	if len(originalTrackingUrl) > 0 {
		this.Ctx.Redirect(302, originalTrackingUrl)
	}
	//else {
	//	adResponse.StatusCode = lib.ERROR_AD_EXPIRED
	//	this.Data["json"] = &adResponse
	//	this.ServeJson()
	//}
}
