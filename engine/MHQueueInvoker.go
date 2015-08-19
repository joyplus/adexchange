package engine

import (
	"adexchange/lib"
	m "adexchange/models"
	"bytes"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/vmihailenco/msgpack.v2"
	"time"
)

func invokeMHQueue(demand *Demand) {
	//timeoutChan := make(chan bool, 1)

	beego.Debug("Start Invoke MHQueue,did:" + demand.AdRequest.Did)

	var adResponse *m.AdResponse
	queueName := generateQueueName(demand)
	//queueChan := make(chan *m.AdResponse)

	go processDemand(demand, queueName)
	go processDemand(demand, queueName)
	//go processAdResponseQueue(queueName, queueChan)
	//go waitQueue(demand.Timeout, timeoutChan)

	//select {
	//case adResponse = <-queueChan:
	//	beego.Debug("Queue Return")
	//	break
	//case <-timeoutChan:
	//	beego.Debug("Queue Timeout")
	//	adResponse = generateErrorResponse(demand.AdRequest, demand.AdspaceKey, lib.ERROR_TIMEOUT_ERROR)
	//	break
	//}
	//beego.Debug("End queue")
	adResponse = getAdResponseFromQueue(queueName)
	if adResponse == nil {
		adResponse = generateErrorResponse(demand.AdRequest, demand.AdspaceKey, lib.ERROR_NO_AD_FROM_QUEUE)
	}

	demand.Result <- adResponse
}

func generateQueueName(demand *Demand) (queueName string) {
	var buffer bytes.Buffer

	strToday := time.Now().Format("2006-01-02")

	buffer.WriteString("_")
	buffer.WriteString("MHQUEUE")
	buffer.WriteString("_")
	buffer.WriteString(strToday)
	buffer.WriteString("_")
	buffer.WriteString(demand.AdRequest.AdspaceKey)
	buffer.WriteString("_")
	buffer.WriteString(demand.AdspaceKey)
	if len(demand.TargetingCode) > 0 {
		buffer.WriteString("_")
		buffer.WriteString(demand.TargetingCode)
	}

	queueName = buffer.String()

	return
}

func processDemand(demand *Demand, queueName string) {

	newDemand := new(Demand)
	newDemand.URL = demand.URL
	newDemand.Timeout = demand.Timeout * 10
	newDemand.AdRequest = demand.AdRequest
	newDemand.AdspaceKey = demand.AdspaceKey
	newDemand.AdSecretKey = demand.AdSecretKey
	newDemand.TargetingCode = demand.TargetingCode
	newDemand.Result = make(chan *m.AdResponse)
	newDemand.Priority = demand.Priority

	//Generate new did for queue invoker

	newDemand.Did = lib.GenerateBid(newDemand.AdRequest.AdspaceKey)
	//beego.Debug(newDemand.Did)

	go invokeMH(newDemand)
	adResponse := <-newDemand.Result

	//sent in MHInvoker
	//go SendDemandLog(adResponse)

	if adResponse.StatusCode == lib.STATUS_SUCCESS {
		go SendDemandResponse(adResponse, queueName)
	}

}

func SendDemandResponse(adResponse *m.AdResponse, queueName string) {

	c := lib.Pool.Get()
	b, err := msgpack.Marshal(adResponse)

	beego.Debug("==========" + beego.AppConfig.String("runmode") + queueName)

	if err == nil {
		c = lib.Pool.Get()
		c.Do("lpush", beego.AppConfig.String("runmode")+queueName, b)
	} else {
		beego.Error(err.Error())
	}
}

//func waitQueue(timeout int, timeoutChan chan bool) {
//	time.Sleep(time.Millisecond * time.Duration(10))
//	timeoutChan <- true
//}

//func processAdResponseQueue(queueName string, queueChan chan *m.AdResponse) {

//	c := lib.Pool.Get()

//	beego.Debug("==========" + beego.AppConfig.String("runmode") + queueName)
//	reply, err := c.Do("rpop", beego.AppConfig.String("runmode")+queueName)

//	if err != nil {
//		beego.Error(err.Error())
//	}
//	var adResponse *m.AdResponse
//	switch reply := reply.(type) {
//	case []byte:
//		b, _ := redis.Bytes(reply, nil)
//		adResponse = getAdResponse(b)
//		break
//	case nil:

//		beego.Info("AdResponse Queue Connection timeout")
//		break
//	default:
//		beego.Info("AdResponse Queue Unknow reply:")
//		beego.Info(reply)
//		break
//	}
//	defer c.Close()

//	queueChan <- adResponse
//	//if queueChan != nil {
//	//	queueChan <- adResponse
//	//}

//}

func getAdResponse(b []byte) (adResponse *m.AdResponse) {
	adResponse = new(m.AdResponse)
	err := msgpack.Unmarshal(b, adResponse)
	if err != nil {
		beego.Critical(err.Error())
	}

	return adResponse
}

func getAdResponseFromQueue(queueName string) (adResponse *m.AdResponse) {

	c := lib.Pool.Get()

	reply, err := c.Do("rpop", beego.AppConfig.String("runmode")+queueName)

	if err != nil {
		beego.Error(err.Error())
	}
	//var adResponse *m.AdResponse
	switch reply := reply.(type) {
	case []byte:
		b, _ := redis.Bytes(reply, nil)
		adResponse = getAdResponse(b)
		break
	case nil:

		beego.Info("AdResponse Queue Connection timeout")
		break
	default:
		beego.Info("AdResponse Queue Unknow reply:")
		beego.Info(reply)
		break
	}
	defer c.Close()

	return
	//if queueChan != nil {
	//	queueChan <- adResponse
	//}

}
