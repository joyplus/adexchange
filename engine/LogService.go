package engine

import (
	"adexchange/lib"
	m "adexchange/models"
	"github.com/astaxie/beego"
	"gopkg.in/vmihailenco/msgpack.v2"
)

var _demandLogPool chan *m.AdResponse
var _mhQueuePool chan *MHQueueData
var _reqLogPool chan *m.AdRequest
var _impLogPool chan *m.AdRequest
var _clkLogPool chan *m.AdRequest

func init() {

	_demandLogPool = make(chan *m.AdResponse, 5000)
	_mhQueuePool = make(chan *MHQueueData, 5000)
	_reqLogPool = make(chan *m.AdRequest, 2000)
	_impLogPool = make(chan *m.AdRequest, 1000)
	_clkLogPool = make(chan *m.AdRequest, 1000)

}

func StartDemandLogService() {

	for {
		adResponse := <-_demandLogPool
		b, err := msgpack.Marshal(adResponse)

		if err == nil {
			c := lib.GetQueuePool().Get()
			c.Do("lpush", beego.AppConfig.String("runmode")+"_ADMUX_DEMAND", b)
			defer c.Close()

		} else {
			beego.Critical(err.Error())
		}
	}

}

func StartMHQueueService() {

	for {
		mhQueueData := <-_mhQueuePool
		b, err := msgpack.Marshal(mhQueueData.AdResponse)

		if err == nil {
			c := lib.GetQueuePool().Get()
			c.Do("lpush", beego.AppConfig.String("runmode")+mhQueueData.QueueName, b)
			defer c.Close()

		} else {
			beego.Critical(err.Error())
		}
	}

}

func StartReqLogService() {

	for {
		adRequest := <-_reqLogPool
		b, err := msgpack.Marshal(adRequest)

		if err == nil {
			c := lib.GetQueuePool().Get()

			c.Do("lpush", beego.AppConfig.String("runmode")+"_ADMUX_REQ", b)
			defer c.Close()

		} else {
			beego.Critical(err.Error())
		}
	}

}

func StartImpLogService() {

	for {
		adRequest := <-_impLogPool
		b, err := msgpack.Marshal(adRequest)

		if err == nil {
			c := lib.GetQueuePool().Get()
			c.Do("lpush", beego.AppConfig.String("runmode")+"_ADMUX_IMP", b)
			defer c.Close()

		} else {
			beego.Critical(err.Error())
		}
	}

}

func StartClkLogService() {

	for {
		adRequest := <-_clkLogPool
		b, err := msgpack.Marshal(adRequest)

		if err == nil {
			c := lib.GetQueuePool().Get()
			c.Do("lpush", beego.AppConfig.String("runmode")+"_ADMUX_CLK", b)
			defer c.Close()

		} else {
			beego.Error(err.Error())
		}
	}

}

func SendDemandLog(adResponse *m.AdResponse) {
	if adResponse != nil {
		_demandLogPool <- adResponse
	}
	//c := lib.Pool.Get()
	//b, err := msgpack.Marshal(adResponse)

	//if err == nil {
	//	c = lib.Pool.Get()
	//	c.Do("lpush", beego.AppConfig.String("runmode")+"_ADMUX_DEMAND", b)
	//} else {
	//	beego.Error(err.Error())
	//}
}

func SendRequestLog(adRequest *m.AdRequest, logType int) {
	if adRequest != nil {
		if logType == 1 {
			_reqLogPool <- adRequest
		} else if logType == 2 {
			_impLogPool <- adRequest
		} else if logType == 3 {
			_clkLogPool <- adRequest
		} else {
			beego.Critical("logtype is wrong!")
		}
	} else {
		beego.Critical("adRequest is null")
	}

}

func SendMHQueue(adResponse *m.AdResponse, queueName string) {
	if adResponse != nil && len(queueName) > 0 {
		mhQueueData := new(MHQueueData)
		mhQueueData.AdResponse = adResponse
		mhQueueData.QueueName = queueName
		_mhQueuePool <- mhQueueData
	}
}
