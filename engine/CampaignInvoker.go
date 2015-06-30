package engine

import (
	"adexchange/lib"
	m "adexchange/models"
	"github.com/astaxie/beego"
	"time"
)

func invokeCampaign(demand *Demand) {

	beego.Debug("Start Invoke Campaign")
	adRequest := demand.AdRequest
	campaigns, err := m.GetCampaigns(adRequest.AdspaceKey, time.Now().Format("2006-01-02"))
	if err != nil {
		beego.Error(err.Error)
	}
	adResponse := new(m.AdResponse)
	adResponse.Bid = adRequest.Bid
	adResponse.SetDemandAdspaceKey(demand.AdspaceKey)

	if len(campaigns) == 0 {

		adResponse.StatusCode = lib.ERROR_NOAD
		demand.Result <- adResponse
	} else {
		random := lib.GetRandomNumber(0, len(campaigns))
		adResponse = mapCampaign(campaigns[random])
		adResponse.Bid = adRequest.Bid
		adResponse.SetDemandAdspaceKey(demand.AdspaceKey)

		demand.Result <- adResponse
	}

}

func mapCampaign(campaign *m.PmpCampaign) (adResponse *m.AdResponse) {

	adResponse = new(m.AdResponse)
	adResponse.StatusCode = lib.STATUS_SUCCESS
	adResponse.SetResponseTime(time.Now().Unix())

	adUnit := new(m.AdUnit)
	adResponse.Adunit = adUnit
	adUnit.Cid = lib.ConvertIntToString(campaign.Id)
	adUnit.ClickUrl = campaign.LandingUrl
	adUnit.ImageUrls = []string{campaign.CreativeUrl}
	adUnit.AdWidth = campaign.Width
	adUnit.AdHeight = campaign.Height

	return adResponse
}
