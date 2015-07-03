package engine

import (
	"adexchange/lib"
	m "adexchange/models"
	"github.com/astaxie/beego"
	"gopkg.in/vmihailenco/msgpack.v2"
	"time"
)

func invokeCampaign(demand *Demand) {

	adRequest := demand.AdRequest
	beego.Debug("Start Invoke Campaign,bid:" + adRequest.Bid)

	adResponse := getCachedAdResponse(adRequest)

	if adResponse == nil {
		adResponse := new(m.AdResponse)
		adResponse.Bid = adRequest.Bid
		adResponse.SetDemandAdspaceKey(demand.AdspaceKey)
		adResponse.SetResponseTime(time.Now().Unix())
		campaigns, err := m.GetCampaigns(adRequest.AdspaceKey, time.Now().Format("2006-01-02"))
		if err != nil {
			beego.Error(err.Error)
		}

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
		setCachedAdResponse(generateCacheKey(adRequest), adResponse)

	} else {

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
	adUnit.CreativeUrls = []string{campaign.CreativeUrl}
	adUnit.AdWidth = campaign.Width
	adUnit.AdHeight = campaign.Height

	return adResponse
}

func generateCacheKey(adRequest *m.AdRequest) string {
	return beego.AppConfig.String("runmode") + "_CAMPAIGN_" + adRequest.AdspaceKey
}

func setCachedAdResponse(cacheKey string, adResponse *m.AdResponse) {
	c := lib.Pool.Get()
	val, err := msgpack.Marshal(adResponse)

	if _, err = c.Do("SET", cacheKey, val); err != nil {
		beego.Error(err.Error())
	}

	_, err = c.Do("EXPIRE", cacheKey, 60)
	if err != nil {
		beego.Error(err.Error())
	}

}

func getCachedAdResponse(adRequest *m.AdRequest) (adResponse *m.AdResponse) {
	c := lib.Pool.Get()

	key := generateCacheKey(adRequest)
	v, err := c.Do("GET", key)
	if err != nil {
		beego.Error(err.Error())
		return nil
	}

	if v == nil {
		return
	}

	adResponse = new(m.AdResponse)
	switch t := v.(type) {
	case []byte:
		err = msgpack.Unmarshal(t, adResponse)
	default:
		err = msgpack.Unmarshal(t.([]byte), adResponse)
	}

	if err != nil {
		beego.Error(err.Error())
	}
	return
}
