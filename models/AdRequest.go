package models

type AdRequest struct {
	Bid              string `form:"bid"`
	AdspaceKey       string `form:"adspaceid"`
	DemandAdspaceKey string `form:"demandadspaceid"`
	AdType           string `form:"adtype"`
	Pkgname          string `form:"pkgname"`
	Appname          string `form:"appname"`
	Conn             string `form:"conn"`
	Carrier          string `form:"carrier"`
	ApiType          string `form:"apitype"`
	Os               string `form:"os"`
	Osv              string `form:"osv"`
	Imei             string `form:"imei"`
	Wma              string `form:"wma"`
	Aid              string `form:"aid"`
	Aaid             string `form:"aaid"`
	Idfa             string `form:"idfa"`
	Oid              string `form:"oid"`
	Uid              string `form:"uid"`
	Device           string `form:"device"`
	Ua               string `form:"ua"`
	Ip               string `form:"ip"`
	Width            string `form:"width"`
	Height           string `form:"height"`
	Pcat             string `form:"pcat"`
	Density          string `form:"density"`
	Lon              string `form:"lon"`
	Lat              string `form:"lat"`
}
