package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func GetCampaigns(demandAdspaceKey string, adDate string) (campaigns []*PmpCampaign, err error) {
	o := orm.NewOrm()

	sql := "select * from pmp_campaign as campaign inner join pmp_demand_adspace as demand on campaign.demand_adspace_id=demand.id where demand.demand_adspace_key=? and campaign.end_date>=? and ?>=campaign.start_date and campaign_status=1"

	paramList := []interface{}{demandAdspaceKey, adDate, adDate}

	_, err = o.Raw(sql, paramList).QueryRows(&campaigns)

	if err != nil {
		beego.Critical(err.Error())
	}

	return campaigns, err

}
