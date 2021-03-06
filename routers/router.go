// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"adexchange/controllers"
	"github.com/astaxie/beego"
)

//func init() {
//	ns := beego.NewNamespace("/v1",
//		beego.NSNamespace("/object",
//			beego.NSInclude(
//				&controllers.ObjectController{},
//			),
//		),
//		beego.NSNamespace("/user",
//			beego.NSInclude(
//				&controllers.UserController{},
//			),
//		),
//	)
//	beego.AddNamespace(ns)
//}

func init() {

	beego.Info("admux start")
	beego.Router("/api/request", &controllers.RequestController{}, "*:RequestAd")
	beego.Router("/api/trackimp", &controllers.RequestController{}, "*:TrackImp")
	beego.Router("/api/trackclk", &controllers.RequestController{}, "*:TrackClk")
	beego.Router("/api/viewad", &controllers.ViewAdController{}, "*:ViewAd")
	beego.Router("/api/clientreq", &controllers.ClientRequestController{}, "*:RequestAd4Client")
	beego.Router("/api/webad", &controllers.WebviewRequestController{}, "*:WebviewReq")

}
