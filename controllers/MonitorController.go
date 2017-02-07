package controllers

import (
	m "adexchange/models"

	"github.com/astaxie/beego"
)

type MonitorController struct {
	beego.Controller
}

//Request Ad
func (this *MonitorController) UpdateStatus() {

	adResponse := new(m.AdResponse)
	beego.Debug("Enter update status")

	this.Data["json"] = &adResponse
	this.ServeJSON()

}
