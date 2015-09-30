package controllers

import (
	"adexchange/engine"
	"adexchange/lib"
	m "adexchange/models"
	"encoding/json"
	"github.com/astaxie/beego"
)

type LTVS2SController struct {
	beego.Controller
}

//Handle active request
func (this *LTVS2SController) handleActiveReq() {

	flg := true
	var eventRequest m.EventRequest
	response := new(m.BaseResponse)
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &eventRequest)
	if err != nil {
		beego.Critical(err.Error())
		response.StatusCode = lib.ERROR_JSON_UNMARSHAL_FAILED
		flg = false
	} else {
		response.StatusCode = lib.STATUS_SUCCESS
	}

	if flg {
		beego.Debug("todo")
	}

	engine.SendEventRequestLog(&eventRequest)
	this.Data["json"] = response
	this.ServeJson()
}
