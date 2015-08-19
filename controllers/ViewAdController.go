package controllers

//import (
//	"github.com/astaxie/beego"
//)

type ViewAdController struct {
	BaseController
}

func (this *ViewAdController) ViewAd() {

	cacheKey := this.GetString("id")
	adResponse := GetCachedAdResponse(cacheKey)

	//beego.Debug(adResponse)

	if adResponse != nil {
		this.TplNames = "tpl/ad.html"

		adParam := map[string][]string{"clkTrackingUrls": adResponse.Adunit.ClkTrackingUrls, "implTrackingUrls": adResponse.Adunit.ImpTrackingUrls, "imgUrls": adResponse.Adunit.CreativeUrls}

		this.Data["AD"] = adParam
		this.Data["clickUrl"] = adResponse.Adunit.ClickUrl
		this.Data["width"] = adResponse.Adunit.AdWidth
		this.Data["height"] = adResponse.Adunit.AdHeight
		this.Render()
	}

}
