package engine

import (
	"adexchange/lib"
	m "adexchange/models"
	"github.com/astaxie/beego"
	"gopkg.in/vmihailenco/msgpack.v2"
)

var _demandLogPool chan *m.AdResponse

func init() {

	_demandLogPool = make(chan *m.AdResponse, 100)

}

func StartDemandLogService() {

	c := lib.Pool.Get()

	for {
		adResponse := <-_demandLogPool
		b, err := msgpack.Marshal(adResponse)

		if err == nil {
			c = lib.Pool.Get()
			c.Do("lpush", "ADMUX_DEMAND", b)
		} else {
			beego.Error(err.Error())
		}
	}

	defer c.Close()
}

func SendDemandLog(adResponse *m.AdResponse) {
	//if adResponse != nil {
	//	_demandLogPool <- adResponse
	//}
	c := lib.Pool.Get()
	b, err := msgpack.Marshal(adResponse)

	if err == nil {
		c = lib.Pool.Get()
		c.Do("lpush", "ADMUX_DEMAND", b)
	} else {
		beego.Error(err.Error())
	}
}
