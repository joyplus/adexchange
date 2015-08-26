package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func GetCampaigns(demandAdspaceKey string, adDate string, width int, height int) (campaigns []*PmpCampaignCreative, err error) {
	o := orm.NewOrm()

	sql := "select creative.name, creative.width,creative.height,creative.creative_url,creative.landing_url,creative.imp_tracking_url from pmp_campaign as campaign inner join pmp_demand_adspace as demand on campaign.demand_adspace_id=demand.id inner join pmp_campaign_creative as creative on campaign.id = creative.campaign_id where demand.demand_adspace_key=? and campaign.end_date>=? and ?>=campaign.start_date and campaign_status=1 and creative.width=? and creative.height=? "

	paramList := []interface{}{demandAdspaceKey, adDate, adDate, width, height}

	_, err = o.Raw(sql, paramList).QueryRows(&campaigns)

	if err != nil {
		beego.Critical(err.Error())
	}

	return campaigns, err

}
