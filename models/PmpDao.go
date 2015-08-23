package models

import (
	"adexchange/lib"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func GetMatrixData() (adspaceMap map[string]AdspaceData, adspaceDemandMap map[string][]string, err error) {
	o := orm.NewOrm()

	sql := "select matrix.priority,matrix.pmp_adspace_id, adspace.pmp_adspace_key, matrix.demand_id as demand_id,demand.demand_adspace_key as demand_adspace_key,demand.secret_key as demand_secret_key, app.pkg_name,app.app_name,app.pcat,app.ua from pmp_adspace_matrix as matrix inner join pmp_adspace as adspace on matrix.pmp_adspace_id=adspace.id inner join pmp_demand_adspace as demand on matrix.demand_adspace_id=demand.id left join pmp_app_info as app on app.id=demand.app_id order by adspace.pmp_adspace_key,matrix.priority desc"

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
	adspaceDemandMap = make(map[string][]string)

	for _, record := range dataList {
		adspaceData := AdspaceData{AdspaceKey: record.DemandAdspaceKey}
		adspaceData.DemandId = record.DemandId
		adspaceData.SecretKey = record.DemandSecretKey
		adspaceData.Priority = record.Priority
		adspaceData.AppName = record.AppName
		adspaceData.PkgName = record.PkgName
		adspaceData.Pcat = record.Pcat
		adspaceData.Ua = record.Ua

		adspaceMap[record.PmpAdspaceKey+"_"+record.DemandAdspaceKey] = adspaceData

		if oldAdspaceKey != record.PmpAdspaceKey {
			oldAdspaceKey = record.PmpAdspaceKey

			if pmpDemandInfo != nil {
				aryDemandAdspaceKey := pmpDemandInfo.GetDemandAdspaceKeys()
				adspaceDemandMap[pmpDemandInfo.AdspaceKey] = aryDemandAdspaceKey
			}
			pmpDemandInfo = new(PmpDemandInfo)
			pmpDemandInfo.InitDemandAdspace()
			pmpDemandInfo.AdspaceKey = record.PmpAdspaceKey
			pmpDemandInfo.AddDemandAdspace(record.DemandAdspaceKey)
		} else {
			pmpDemandInfo.AddDemandAdspace(record.DemandAdspaceKey)
		}
	}

	aryDemandAdsapceKey := pmpDemandInfo.GetDemandAdspaceKeys()
	adspaceDemandMap[pmpDemandInfo.AdspaceKey] = aryDemandAdsapceKey

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

	return demandMap, nil
}

//todo
func validUrl(url string) bool {
	return true
}

//adDate: 2006-01-02
func GetAvbDemandMap(adDate string) (avbDemandMap map[string]*AvbDemand, err error) {
	o := orm.NewOrm()

	beego.Debug("Get Avb Demand Map")

	var records []*AvbDemand
	sql := "select allocation.id as allocation_id,allocation.pmp_adspace_id, allocation.demand_adspace_id, pmp.pmp_adspace_key,demand.demand_adspace_key,allocation.imp as plan_imp, allocation.clk as plan_clk,report.imp as actual_imp, report.clk as actual_clk from pmp_daily_allocation as allocation left join pmp_daily_report as report on allocation.ad_date=report.ad_date and allocation.pmp_adspace_id=report.pmp_adspace_id and allocation.demand_adspace_id=report.demand_adspace_id inner join pmp_adspace as pmp on allocation.pmp_adspace_id=pmp.id inner join pmp_demand_adspace as demand on allocation.demand_adspace_id=demand.id where allocation.ad_date=?"

	paramList := []interface{}{adDate}

	_, err = o.Raw(sql, paramList).QueryRows(&records)

	if err != nil {
		beego.Critical(err.Error())
		return
	}

	sql = "select targeting_type, targeting_code, plan_imp, plan_clk, actual_imp, actual_clk from pmp_daily_allocation_detail where allocation_id=? and targeting_type in ('PROVINCE','CITY') "

	avbDemandMap = make(map[string]*AvbDemand)

	for _, record := range records {
		var detailList []*AllocationDetail
		paramList = []interface{}{record.AllocationId}
		_, err = o.Raw(sql, paramList).QueryRows(&detailList)
		if err != nil {
			beego.Critical(err.Error())
			continue
		}

		if detailList != nil && len(detailList) > 0 {
			for _, detail := range detailList {
				record.SetDetailAllocation(detail)
			}
		}

		avbDemandMap[record.PmpAdspaceKey+"_"+record.DemandAdspaceKey] = record
		//if record.PlanImp > record.ActualImp {
		//	avbDemandMap[record.PmpAdspaceKey+"_"+record.DemandAdspaceKey] = true
		//}
	}

	return avbDemandMap, err
}

func GetPmpInfo() (pmpAdspaceMap map[string]PmpInfo, err error) {
	o := orm.NewOrm()

	sql := "select pmp_adspace_key, creative_type, tpl_name from pmp_adspace where status=0"

	var dataList []PmpInfo

	_, err = o.Raw(sql).QueryRows(&dataList)

	if err != nil {
		beego.Critical(err.Error())
		return
	}

	pmpAdspaceMap = make(map[string]PmpInfo)

	for _, record := range dataList {
		pmpAdspaceMap[record.PmpAdspaceKey] = record
	}

	return
}

func GetTplSet() (tplHashSet *lib.HashSet, err error) {
	o := orm.NewOrm()

	sql := "select distinct(tpl_name) as tpl_name from pmp_adspace "

	var dataList []PmpInfo

	_, err = o.Raw(sql).QueryRows(&dataList)

	if err != nil {
		beego.Critical(err.Error())
		return
	}

	tplHashSet = lib.NewHashSet()

	for _, record := range dataList {
		tplHashSet.Add(record.TplName)
	}

	return
}
