package controllers

//import (
//	"adexchange/lib"
//	m "adexchange/models"
//	"github.com/astaxie/beego"
//)

type ViewAdController struct {
	BaseController
}

func (this *RequestController) ViewAd() {
	
	cacheKey := this.GetString("id")
	adResponse := GetCachedAdResponse(cacheKey);
	
	this.TplNames = "tpl/ad.html"
	
	adParam :=map[string][]string{"clkTrackingUrls": adResponse.Adunit.ClkTrackingUrls, "implTrackingUrls": adResponse.Adunit.ImpTrackingUrls, "imgUrls": adResponse.Adunit.CreativeUrls}

	this.Data["AD"] = adParam
	this.Data["width"] = adResponse.Adunit.AdWidth
	this.Data["height"] = adResponse.Adunit.AdHeight
	this.Render()
}
