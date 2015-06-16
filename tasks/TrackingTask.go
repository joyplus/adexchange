package tasks

import (
	"adexchange/lib"
	m "adexchange/models"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/vmihailenco/msgpack.v2"
)

func HandleImp() {

	c := lib.Pool.Get()
	var adRequest m.AdRequest
	for {
		b, _ := redis.Bytes(c.Do("brpop", "ADMUX_IMP", "0"))

		err := msgpack.Unmarshal(b, &adRequest)

		if err != nil {
			beego.Error(err.Error())
		} else {
			beego.Debug(adRequest)
		}

	}

	defer c.Close()
}
