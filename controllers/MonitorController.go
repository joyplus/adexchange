package controllers

import (
	"adexchange/engine"
	"adexchange/lib"
	m "adexchange/models"
	"github.com/astaxie/beego"
)

type MonitorController struct {
	beego.Controller
}

//Request Ad
func (this *MonitorController) UpdateStatus() {

	adResponse := new(m.AdResponse)
	beego.Debug("Enter Request ad")

	m.InitEngineData()

	this.Data["json"] = &adResponse
	this.ServeJson()

}
