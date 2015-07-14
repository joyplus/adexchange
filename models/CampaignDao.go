package models

import (
	"github.com/astaxie/beego/orm"
)

func GetCampaigns(pmpAdspaceKey string, adDate string) (campaigns []*PmpCampaign, err error) {
	o := orm.NewOrm()

	sql := "select * from pmp_campaign as campaign inner join pmp_campaign_matrix as matrix on campaign.id=matrix.pmp_campaign_id inner join pmp_adspace as pmp on matrix.pmp_adspace_id=pmp.id where pmp.pmp_adspace_key=? and campaign.end_date>=? and ?>=campaign.start_date and campaign.status=1"

	paramList := []interface{}{pmpAdspaceKey, adDate, adDate}

	_, err = o.Raw(sql, paramList).QueryRows(&campaigns)

	return campaigns, err

}
