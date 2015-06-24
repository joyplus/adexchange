package controllers

import (
	"adexchange/lib"
	m "adexchange/models"
	"github.com/astaxie/beego"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type BaseController struct {
	beego.Controller
}

//Send log to the queue
func SendLog(adRequest m.AdRequest, logType int) {

	b, err := msgpack.Marshal(adRequest)

	if err == nil {
		c := lib.Pool.Get()
		c.Do("lpush", getQueueName(logType), b)

		defer c.Close()
	} else {

		beego.Error(err.Error())
	}

}

func getQueueName(logType int) string {
	if logType == 1 {
		return "ADMUX_REQ"
	} else if logType == 2 {
		return "ADMUX_IMP"
	} else if logType == 3 {
		return "ADMUX_CLK"
	} else {
		return ""
	}
}
