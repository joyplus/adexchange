package models

import (
	"adexchange/lib"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func GetMatrixData() (adspaceMap map[string]AdspaceData, adspaceDemandMap map[string][]int, err error) {
	o := orm.NewOrm()

	sql := "select matrix.pmp_adspace_id, adspace.pmp_adspace_key, matrix.demand_id as demand_id,demand.demand_adspace_key as demand_adspace_key,demand.secret_key as demand_adspace_secret from pmp_adspace_matrix as matrix inner join pmp_adspace as adspace on matrix.pmp_adspace_id=adspace.id inner join pmp_demand_adspace as demand on matrix.demand_adspace_id=demand.id order by adspace.pmp_adspace_key,matrix.priority"

	var dataList []PmpAdplaceInfo

	_, err = o.Raw(sql).QueryRows(&dataList)

	if err != nil {
		return nil, nil, err
	}

	var oldAdspaceKey string
	var pmpDemandInfo *PmpDemandInfo

	//key:<adspace_key>_<demand_id>; value:<demand_adspace_key>,<demand_secret_key>
	adspaceMap = make(map[string]AdspaceData)

	//key:<adspace_key>; value:<demand_id1>,<demand_id2>...
	adspaceDemandMap = make(map[string][]int)

	for _, record := range dataList {
		adspaceData := AdspaceData{AdspaceKey: record.DemandAdspaceKey}
		adspaceData.SecretKey = record.DemandSecretKey
		adspaceMap[record.PmpAdspaceKey+"_"+lib.ConvertIntToString(record.DemandId)] = adspaceData

		if oldAdspaceKey != record.PmpAdspaceKey {
			oldAdspaceKey = record.PmpAdspaceKey

			if pmpDemandInfo != nil {
				demandIds := pmpDemandInfo.GetDemandIds()
				adspaceDemandMap[pmpDemandInfo.AdspaceKey] = demandIds
			}
			pmpDemandInfo = new(PmpDemandInfo)
			pmpDemandInfo.InitDemand()
			pmpDemandInfo.AdspaceKey = record.PmpAdspaceKey
			pmpDemandInfo.AddDemand(record.DemandId)
		} else {
			pmpDemandInfo.AddDemand(record.DemandId)
		}
	}

	demandIds := pmpDemandInfo.GetDemandIds()
	adspaceDemandMap[pmpDemandInfo.AdspaceKey] = demandIds

	return adspaceMap, adspaceDemandMap, err
}

func GetDemandInfo() (demandMap map[int]DemandInfo, err error) {
	o := orm.NewOrm()

	sql := "select id as demand_id, request_url_template, name, timeout, invoke_func_name from pmp_demand_platform_desk"

	var dataList []DemandInfo

	_, err = o.Raw(sql).QueryRows(&dataList)

	if err != nil {
		return nil, err
	}

	demandMap = make(map[int]DemandInfo)

	for _, record := range dataList {
		demandMap[record.DemandId] = record
	}

	beego.Debug(demandMap)

	return demandMap, nil
}

//todo
func validUrl(url string) bool {
	return true
}

//adDate: 2006-01-02
func GetAvbDemandMap(adDate string) (avbDemandMap map[string]bool, err error) {
	o := orm.NewOrm()

	beego.Debug("Start update demand daily report")

	var records []*AvbDemand
	sql := "select pmp.pmp_adspace_key,demand.demand_adspace_key from pmp_daily_allocation as allocation inner join pmp_daily_report as report on allocation.ad_date=report.ad_date and allocation.pmp_adspace_id=report.pmp_adspace_id and allocation.demand_adspace_id=report.demand_adspace_id and allocation.imp>report.imp inner join pmp_adspace as pmp on allocation.pmp_adspace_id=pmp.id inner join pmp_demand_adspace as demand on allocation.demand_adspace_id=demand.id where allocation.ad_date=?"

	paramList := []interface{}{adDate}

	_, err = o.Raw(sql, paramList).QueryRows(&records)

	if err != nil {
		return avbDemandMap, err
	}

	avbDemandMap = make(map[string]bool)

	for _, record := range records {
		avbDemandMap[record.PmpAdspaceKey+"_"+record.DemandAdspaceKey] = true
	}

	return avbDemandMap, err
}
