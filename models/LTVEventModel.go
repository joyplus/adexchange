package models

//import (
//	"time"
//)

type EventRequest struct {
	AppKey     string
	Event      string
	DeviceList []*Device
}

type Device struct {
	OS   int
	Idfa string
	Mac  string
	Imei string
	Aid  string
}

type BaseResponse struct {
	StatusCode int
	ErrorMsg   string
}

type InstallRequest struct {
	Spreadurl  string `form:"spreadurl"`
	Spreadname string `form:"spreadname"`
	Clicktime  int    `form:"clicktime"`
	Ua         string `form:"ua"`
	Uip        string `form:"uip"`
	Appkey     string `form:"appkey"`
	Activetime int    `form:"activetime"`
	Osversion  string `form:"osversion"`
	Devicetype int    `form:"devicetype"`
	Idfa       string `form:"idfa"`
	Mac        string `form:"mac"`
	Gpid       string `form:"gpid"`
	Aid        string `form:"aid"`
	Skey       string `form:"skey"`
}
