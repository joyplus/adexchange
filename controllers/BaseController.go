package controllers

import (
	"adexchange/lib"
	m "adexchange/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/vmihailenco/msgpack.v2"
	//"strings"
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
	prefix := beego.AppConfig.String("runmode") + "_"

	if logType == 1 {
		return prefix + "ADMUX_REQ"
	} else if logType == 2 {
		return prefix + "ADMUX_IMP"
	} else if logType == 3 {
		return prefix + "ADMUX_CLK"
	} else {
		return ""
	}
}

func GetClientIP(input *context.BeegoInput) string {
	//ips := input.Proxy()
	//if len(ips) > 0 && ips[0] != "" {
	//	return ips[0]
	//}
	//ip := strings.Split(input.Request.RemoteAddr, ":")
	//if len(ip) > 0 {
	//	return ip[0]
	//}
	return input.IP()
}

func SetCachedClkUrl(cacheKey string, clkUrl string) (err error) {
	c := lib.Pool.Get()
	prefix := beego.AppConfig.String("runmode") + "_"

	if _, err = c.Do("SET", prefix+cacheKey, clkUrl); err != nil {
		beego.Error(err.Error())
	}

	_, err = c.Do("EXPIRE", prefix+cacheKey, 300)
	if err != nil {
		beego.Error(err.Error())
	}

	return
}

func GetCachedClkUrl(cacheKey string) (clkUrl string) {
	c := lib.Pool.Get()
	prefix := beego.AppConfig.String("runmode") + "_"
	beego.Debug(prefix + cacheKey)
	clkUrl, err := redis.String(c.Do("GET", prefix+cacheKey))

	if err != nil {
		beego.Error(err.Error())
	}

	return
}

func SetCachedAdResponse(cacheKey string, adResponse *m.AdResponse) {
	c := lib.Pool.Get()
	prefix := beego.AppConfig.String("runmode") + "_"

	val, err := msgpack.Marshal(adResponse)

	if _, err = c.Do("SET", prefix+cacheKey, val); err != nil {
		beego.Error(err.Error())
	}

	_, err = c.Do("EXPIRE", prefix+cacheKey, 120)
	if err != nil {
		beego.Error(err.Error())
	}
}

func GetCachedAdResponse(cacheKey string) (adResponse *m.AdResponse) {
	c := lib.Pool.Get()
	prefix := beego.AppConfig.String("runmode") + "_"

	v, err := c.Do("GET", prefix+cacheKey)
	if err != nil {
		beego.Error(err.Error())
		return nil
	}

	if v == nil {
		return
	}

	adResponse = new(m.AdResponse)
	switch t := v.(type) {
	case []byte:
		err = msgpack.Unmarshal(t, adResponse)
	default:
		err = msgpack.Unmarshal(t.([]byte), adResponse)
	}

	if err != nil {
		beego.Error(err.Error())
	}
	return
}
