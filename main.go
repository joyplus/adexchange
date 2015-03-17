package main

import (
	"github.com/astaxie/beego"
	//"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/ddliu/go-httpclient"
	"strings"
	"sync"
	_ "testweb/routers"
)

var c1, c2 httpclient.HttpClient
var once sync.Once
var log logs.BeeLogger

type Demand struct {
	url    string
	client httpclient.HttpClient
	result chan string
}

const serverUrl1 = "http://ad.madserving.com/adcall_j/"
const serverUrl2 = "http://ad.madserving.com/adcall_j/bidrequest?adspaceid=MAHAD90000001&adtype=2&width=1280&height=200&pid=1&pcat=111&media=1&bid=0000000006&ip=&uid=6ec9fe2b94777ace00d5f67528570cf9&idfa=6ec9fe2b94777ace00d5f67528570cf9&oid=6ec9fe2b94777ace00d5f67528570cf9&vid=6ec9fe2b94777ace00d5f67528570cf9&aid=6ec9fe2b94777ace00d5f67528570cf9&imei=6ec9fe2b94777ace00d5f67528570cf9&aaid=6ec9fe2b94777ace00d5f67528570cf9&wma=6ec9fe2b94777ace00d5f67528570cf9&os=0&osv=4.3&ua=Mozilla/5.0%20(Linux;%20Android%204.3;%20zh-cn;%20ME525+%20Build/)%20AppleWebKit/534.30%20(KHTML,%20like%20Gecko)%20Version/4.0%20Mobile%20Safari/534.30&pkgname=&appname=&conn=0&carrier=0&apitype=0&density=&cell=&device=iPhone5s"

type MainController struct {
	beego.Controller
}

func setup() {

	log := logs.NewLogger(10000)
	log.SetLogger("console", `{"level":5}`)

	//c1 := httpclient.NewHttpClient().Defaults(httpclient.Map{
	//	httpclient.OPT_USERAGENT: "browser1", httpclient.OPT_CONNECTTIMEOUT_MS: 300, httpclient.OPT_TIMEOUT_MS: 80,
	//})

	//c1.Get("http://www.baidu.com/", nil)

	//c2 := httpclient.NewHttpClient().Defaults(httpclient.Map{
	//	httpclient.OPT_USERAGENT: "browser2", httpclient.OPT_CONNECTTIMEOUT_MS: 300, httpclient.OPT_TIMEOUT_MS: 80,
	//})

	//c2.Get("http://www.baidubee.com/", nil)

}

func main() {
	setup()
	beego.Router("/bidrequest", &MainController{})
	beego.Router("/api", &MainController{})

	beego.Run()
}

func (this *MainController) Get() {

	requestString := this.Ctx.Request.RequestURI

	demands := GetDemandUrls(requestString)

	for _, demand := range demands {
		go ProcessDemand(demand)
	}

	var result string

	for _, demand := range demands {

		tmp := <-demand.result
		beego.BeeLogger.Debug("result:" + tmp)

		if len(result) == 0 {
			if strings.Contains(tmp, "imgurl") {
				result = tmp
			}
		}

		beego.BeeLogger.Debug("Finish")
	}

	this.Ctx.WriteString(result)
}

func GetDemandUrls(requestString string) []*Demand {

	beego.BeeLogger.Debug("Request String:" + requestString)
	demand1 := new(Demand)
	demand1.client = c1
	demand1.url = serverUrl1 + requestString
	demand1.result = make(chan string)

	demand2 := new(Demand)
	demand2.client = c2
	demand2.url = serverUrl2 + requestString
	demand2.result = make(chan string)

	demands := []*Demand{demand1, demand2}

	return demands
}

func ProcessDemand(demand *Demand) {

	res, err := demand.client.Begin().Get(demand.url, nil)
	if err != nil {
		beego.BeeLogger.Error("System error:" + err.Error())

	}
	bodyString, _ := res.ToString()

	demand.result <- bodyString
}
