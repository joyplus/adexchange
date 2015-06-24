package controllers

import (
	"adexchange/lib"
	m "adexchange/models"
	"github.com/astaxie/beego"
)

type TrackingController struct {
	beego.Controller
}

func (this *RequestController) TrackImp() {
	adRequest := m.AdRequest{}
	adResponse := new(m.AdResponse)
	beego.Debug("Enter Tracking imp")

	if err := this.ParseForm(&adRequest); err != nil {
		adResponse.StatusCode = lib.ERROR_PARSE_REQUEST
	} else {
		adResponse.StatusCode = lib.STATUS_SUCCESS
		SendLog(adRequest, 2)

	}

	this.Data["json"] = &adResponse
	this.ServeJson()

}

func (this *RequestController) TrackClk() {
	adRequest := m.AdRequest{}
	adResponse := new(m.AdResponse)
	beego.Debug("Enter Tracking clk")

	if err := this.ParseForm(&adRequest); err != nil {
		adResponse.StatusCode = lib.ERROR_PARSE_REQUEST
	} else {
		adResponse.StatusCode = lib.STATUS_SUCCESS
		SendLog(adRequest, 3)
	}

	this.Data["json"] = &adResponse
	this.ServeJson()

}
