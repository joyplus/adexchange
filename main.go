package main

import (
	"adexchange/lib"
	m "adexchange/models"
	_ "adexchange/routers"
	"adexchange/tasks"
	"adexchange/tools"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

//var c1, c2 httpclient.HttpClient
//var once sync.Once

//type Demand struct {
//	url    string
//	client httpclient.HttpClient
//	result chan string
//}

//type LogData struct {
//	domain       string
//	hasAd        bool
//	responseTime int
//}

//const serverUrl1 = "http://ad.madserving.com/adcall_j/"
//const serverUrl2 = "http://ad.madserving.com/adcall_j/bidrequest?adspaceid=MAHAD90000001&adtype=2&width=1280&height=200&pid=1&pcat=111&media=1&bid=0000000006&ip=&uid=6ec9fe2b94777ace00d5f67528570cf9&idfa=6ec9fe2b94777ace00d5f67528570cf9&oid=6ec9fe2b94777ace00d5f67528570cf9&vid=6ec9fe2b94777ace00d5f67528570cf9&aid=6ec9fe2b94777ace00d5f67528570cf9&imei=6ec9fe2b94777ace00d5f67528570cf9&aaid=6ec9fe2b94777ace00d5f67528570cf9&wma=6ec9fe2b94777ace00d5f67528570cf9&os=0&osv=4.3&ua=Mozilla/5.0%20(Linux;%20Android%204.3;%20zh-cn;%20ME525+%20Build/)%20AppleWebKit/534.30%20(KHTML,%20like%20Gecko)%20Version/4.0%20Mobile%20Safari/534.30&pkgname=&appname=&conn=0&carrier=0&apitype=0&density=&cell=&device=iPhone5s"

//const serverUrl1 = "http://ad.madserving.com"
//const serverUrl2 = "http://ad.madserving.com"

//type MainController struct {
//	beego.Controller
//}

//func setup() {
//	//flag.Parse()
//	//pool = newPool(*redisServer, *redisPassword)
//	beego.SetLogger("file", `{"filename":"logs/admux.log"}`)
//	beego.SetLogFuncCall(true)

//	//c1 := httpclient.NewHttpClient().Defaults(httpclient.Map{
//	//	httpclient.OPT_USERAGENT: "browser1", httpclient.OPT_CONNECTTIMEOUT_MS: 300, httpclient.OPT_TIMEOUT_MS: 80,
//	//})

//	//c1.Get("http://www.baidu.com/", nil)

//	//c2 := httpclient.NewHttpClient().Defaults(httpclient.Map{
//	//	httpclient.OPT_USERAGENT: "browser2", httpclient.OPT_CONNECTTIMEOUT_MS: 300, httpclient.OPT_TIMEOUT_MS: 80,
//	//})

//	//c2.Get("http://www.baidubee.com/", nil)

//}

func main() {
	//setup()
	//beego.EnableAdmin = true
	//beego.AdminHttpAddr = "localhost"
	//beego.AdminHttpPort = 8888
	//runtime.GOMAXPROCS(runtime.NumCPU())

	beego.SetLogger("file", `{"filename":"logs/admux.log"}`)
	beego.SetLogFuncCall(true)
	logLevel, _ := beego.AppConfig.Int("log_level")
	beego.SetLevel(logLevel)

	orm.Debug, _ = beego.AppConfig.Bool("orm_debug")
	tools.Init("ip.dat")
	m.Connect()

	lib.Pool = lib.NewPool(beego.AppConfig.String("redis_server"), "")
	tasks.InitEngineData()
	tasks.CheckAvbDemand()
	initDuration, _ := beego.AppConfig.Int("init_duration")
	avbCheckDuration, _ := beego.AppConfig.Int("avb_check_duration")
	go tasks.ScheduleInit(initDuration)
	go tasks.CheckAvbDemandInit(avbCheckDuration)

	//go engine.StartDemandLogService()

	beego.Run()
}

//func (this *MainController) Get() {

//	requestString := this.Ctx.Request.RequestURI

//	demands := GetDemandUrls(requestString)

//	for _, demand := range demands {
//		go ProcessDemand(demand)
//	}

//	var result, tmp string

//	for _, demand := range demands {

//		tmp = <-demand.result
//		//beego.BeeLogger.Debug("result:" + tmp)

//		if len(result) == 0 {
//			if strings.Contains(tmp, "imgurl") {
//				result = tmp
//				beego.Info("success!")
//			}
//		}

//	}
//	if len(result) == 0 {
//		result = tmp
//		beego.Info("no ads")
//	}
//	//jsonStr := `{"MAH3AD90000001": {"adspaceid": "MAH3AD90000001", "returncode": 405} }`
//	//var dat map[string]interface{}
//	//if err := json.Unmarshal([]byte(jsonStr), &dat); err == nil {
//	//	beego.Info("==============json str 转map=======================")
//	//	beego.Info(dat)
//	//	beego.Info(dat["MAH3AD90000001"])
//	//}

//	c := pool.Get()
//	c.Do("lpush", "ADMUX_LOG", "test")

//	defer c.Close()
//	this.Ctx.WriteString(result)
//}

//func GetDemandUrls(requestString string) []*Demand {

//	//beego.BeeLogger.Debug("Request String:" + requestString)
//	demand1 := new(Demand)
//	demand1.client.WithOptions(httpclient.Map{httpclient.OPT_CONNECTTIMEOUT_MS: 200, httpclient.OPT_TIMEOUT_MS: 200})
//	demand1.url = serverUrl1 + requestString
//	demand1.result = make(chan string)

//	demand2 := new(Demand)
//	demand2.client.WithOptions(httpclient.Map{httpclient.OPT_CONNECTTIMEOUT_MS: 200, httpclient.OPT_TIMEOUT_MS: 200})

//	demand2.url = serverUrl2 + requestString
//	demand2.result = make(chan string)

//	demands := []*Demand{demand1, demand2}

//	return demands
//}

//func ProcessDemand(demand *Demand) {

//	var bodyString string
//	res, err := demand.client.Get(demand.url, nil)
//	if err != nil {

//		if httpclient.IsTimeoutError(err) {
//			bodyString = "timeout"
//			beego.Info("timeout")
//		} else {
//			bodyString = "error"
//		}

//	} else {
//		bodyString, _ = res.ToString()
//	}

//	demand.result <- bodyString
//}
