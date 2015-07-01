package controllers

import (
	"adexchange/lib"
	m "adexchange/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"gopkg.in/vmihailenco/msgpack.v2"
	"strings"
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
	prefix := beego.AppConfig.String("runmode")

	if logType == 1 {
		return prefix + "_ADMUX_REQ"
	} else if logType == 2 {
		return prefix + "_ADMUX_IMP"
	} else if logType == 3 {
		return prefix + "_ADMUX_CLK"
	} else {
		return ""
	}
}

func GetClientIP(input *context.BeegoInput) string {
	ips := input.Proxy()
	if len(ips) > 0 && ips[0] != "" {
		return ips[0]
	}
	ip := strings.Split(input.Request.RemoteAddr, ":")
	if len(ip) > 0 {
		return ip[0]
	}
	return ""
}
