// AdTrackingController
package controllers

import (
	"adexchange/engine"
	"adexchange/lib"
	m "adexchange/models"
	"github.com/astaxie/beego"
	"time"
//	"net/url"
)

type AdTrackingController struct {
	BaseController
}

//Request Ad image page
func (this *AdTrackingController) RequestAd() {

	adRequest := m.AdRequest{}
	adResponse := new(m.AdResponse)
	beego.Debug("Enter Request ad")
	if err := this.ParseForm(&adRequest); err != nil {

		adResponse.StatusCode = lib.ERROR_PARSE_REQUEST
	} else {
		adRequest.RequestTime = time.Now().Unix()
		tmp := engine.InvokeDemand(&adRequest)

		if tmp == nil {
			adResponse.StatusCode = lib.ERROR_UNKNON_ERROR
		} else {
			adResponse = tmp
		}
		adRequest.StatusCode = adResponse.StatusCode
		SendLog(adRequest, 1)
		//if err != nil {
		//	beego.Debug("Enter sss ad")
		//	if e, ok := err.(*lib.SysError); ok {
		//		adResponse.StatusCode = e.ErrorCode
		//	} else {
		//		adResponse.StatusCode = lib.ERROR_UNKNON_ERROR
		//	}
		//	beego.Debug("Enter ssaass ad")
		//}
	}
	
	this.TplNames = "tpl/ad.html"
	
	adParam :=map[string][]string{"clkTrackingUrls": adResponse.Adunit.ClkTrackingUrls, "implTrackingUrls": adResponse.Adunit.ImpTrackingUrls, "imgUrls": adResponse.Adunit.ImageUrls}

	this.Data["AD"] = adParam
	this.Data["width"] = adResponse.Adunit.AdWidth
	this.Data["height"] = adResponse.Adunit.AdHeight
	this.Render()
	
//	this.Ctx.Output.Body([]byte("hello world."));
	
	

}
