package controllers

import (
	"adexchange/lib"
	m "adexchange/models"
	"github.com/astaxie/beego"
	"gopkg.in/vmihailenco/msgpack.v2"
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

		b, err := msgpack.Marshal(adRequest)
		beego.Debug(b)
		if err == nil {
			c := lib.Pool.Get()
			c.Do("lpush", "ADMUX_IMP", b)

			defer c.Close()
		} else {
			adResponse.StatusCode = lib.ERROR_MSGPACK_IMP
			beego.Error(err.Error())
		}

	}

	this.Data["json"] = &adResponse
	this.ServeJson()

}
