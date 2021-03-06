package engine

import (
	"adexchange/lib"
	m "adexchange/models"
	"bytes"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/vmihailenco/msgpack.v2"
	"time"
)

type MHQueueData struct {
	QueueName  string
	AdResponse *m.AdResponse
}

func invokeMHQueue(demand *Demand) {

	beego.Debug("Start Invoke MHQueue,did:" + demand.Did)

	var adResponse *m.AdResponse
	queueName := generateQueueName(demand)
	queueChan := make(chan *m.AdResponse)
	timeoutChan := make(chan bool, 1)

	go processDemand(demand, queueName, 0)
	go processDemand(demand, queueName, 1)

	//go processDemand(demand, queueName, i)
	//go processAdResponseQueue(queueName, queueChan)
	go waitQueue(demand.Timeout, timeoutChan)
	go processAdResponseFromQueue(queueName, queueChan)
	select {
	case adResponse = <-queueChan:
		beego.Debug("Queue Return adresponse")
		break
	case <-timeoutChan:
		beego.Debug("Queue Timeout")
		adResponse = generateErrorResponse(demand.AdRequest, demand.AdspaceKey, lib.ERROR_TIMEOUT_ERROR)
		break
	}
	//beego.Debug("End queue")
	//t1 := time.Now().UnixNano()
	//adResponse = getAdResponseFromQueue(queueName)
	//t2 := time.Now().UnixNano()

	//duration := int((t2 - t1) / 1000000)

	//if duration > 100 {
	//	beego.Info(fmt.Sprintf("=====Redis duration=====:%d", duration))
	//}
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
	buffer.WriteString("_")
	buffer.WriteString(lib.ConvertIntToString(demand.AdRequest.Width))
	buffer.WriteString("_")
	buffer.WriteString(lib.ConvertIntToString(demand.AdRequest.Height))

	if len(demand.TargetingCode) > 0 {
		buffer.WriteString("_")
		buffer.WriteString(demand.TargetingCode)
	}

	queueName = buffer.String()

	return
}

func processDemand(demand *Demand, queueName string, index int) {

	newDemand := new(Demand)
	newDemand.URL = demand.URL
	newDemand.Timeout = demand.Timeout * 10
	newDemand.AdRequest = demand.AdRequest
	newDemand.RealAdspaceKey = demand.RealAdspaceKey
	newDemand.AdspaceKey = demand.AdspaceKey
	newDemand.AdSecretKey = demand.AdSecretKey
	newDemand.TargetingCode = demand.TargetingCode
	newDemand.Result = make(chan *m.AdResponse)
	newDemand.Priority = demand.Priority

	//Generate new did for queue invoker

	newDemand.Did = lib.GenerateBid(newDemand.AdRequest.AdspaceKey + lib.ConvertIntToString(index))
	//beego.Debug(newDemand.Did)

	go invokeMH(newDemand)
	adResponse := <-newDemand.Result

	//sent in MHInvoker
	//go SendDemandLog(adResponse)

	if adResponse.StatusCode == lib.STATUS_SUCCESS {
		SendMHQueue(adResponse, queueName)
	}

}

//func SendDemandResponse(adResponse *m.AdResponse, queueName string) {

//	c := lib.Pool.Get()
//	b, err := msgpack.Marshal(adResponse)

//	beego.Debug("==========" + beego.AppConfig.String("runmode") + queueName)

//	if err == nil {
//		c = lib.Pool.Get()
//		c.Do("lpush", beego.AppConfig.String("runmode")+queueName, b)
//	} else {
//		beego.Error(err.Error())
//	}
//}

func waitQueue(timeout int, timeoutChan chan bool) {
	time.Sleep(time.Millisecond * time.Duration(timeout))
	timeoutChan <- true
}

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

func processAdResponseFromQueue(queueName string, queueChan chan *m.AdResponse) {

	t1 := time.Now().UnixNano()

	var adResponse *m.AdResponse

	c := lib.GetQueuePool().Get()

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

		beego.Debug("AdResponse Queue return nil")
		break
	default:
		beego.Debug("AdResponse Queue Unknow reply:")
		beego.Debug(reply)
		break
	}

	defer c.Close()

	t2 := time.Now().UnixNano()
	duration := int((t2 - t1) / 1000000)

	if duration > 100 {
		beego.Info(fmt.Sprintf("=====Redis duration=====:%d", duration))
	}

	queueChan <- adResponse
	//return
	//if queueChan != nil {
	//	queueChan <- adResponse
	//}

}
