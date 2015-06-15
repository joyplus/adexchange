package models

import (
	"adexchange/lib"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func InitEngineData() {

	beego.Debug("Start Init Engine Data")
	adspaceMap, adspaceDemandMap, err := GetMatrixData()

	if err != nil {
		panic(err.Error())
	}
	beego.Debug(adspaceMap)
	beego.Debug(adspaceDemandMap)
	demandMap, err := GetDemandInfo()

	beego.Debug(demandMap)

}

func GetMatrixData() (adspaceMap map[string]AdspaceData, adspaceDemandMap map[string][]int, err error) {
	o := orm.NewOrm()

	sql := "select matrix.pmp_adspace_id, adspace.adspace_key, matrix.demand_id as demand_id,demand.adspace_key as demand_adspace_key,demand.secret_key as demand_adspace_secret from pmp_adspace_matrix as matrix inner join pmp_adspace as adspace on matrix.pmp_adspace_id=adspace.id inner join pmp_demand_adspace as demand on matrix.demand_adspace_id=demand.id order by adspace.adspace_key,matrix.priority"

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
		adspaceMap[record.AdspaceKey+"_"+lib.ConvertIntToString(record.DemandId)] = adspaceData

		if oldAdspaceKey != record.AdspaceKey {
			oldAdspaceKey = record.AdspaceKey

			if pmpDemandInfo != nil {
				demandIds := pmpDemandInfo.GetDemandIds()
				adspaceDemandMap[pmpDemandInfo.AdspaceKey] = demandIds
			}
			pmpDemandInfo = new(PmpDemandInfo)
			pmpDemandInfo.InitDemand()
			pmpDemandInfo.AdspaceKey = record.AdspaceKey
			pmpDemandInfo.AddDemand(record.DemandId)
		} else {
			pmpDemandInfo.AddDemand(record.DemandId)
		}
	}

	demandIds := pmpDemandInfo.GetDemandIds()
	adspaceDemandMap[pmpDemandInfo.AdspaceKey] = demandIds

	return adspaceMap, adspaceDemandMap, err
}

func GetDemandInfo() (demandMap map[int]string, err error) {
	o := orm.NewOrm()

	sql := "select id as demand_id, request_url_template as url from pmp_demand_platform_desk"

	var dataList []DemandInfo

	_, err = o.Raw(sql).QueryRows(&dataList)

	if err != nil {
		return nil, err
	}

	demandMap = make(map[int]string)

	for _, record := range dataList {
		if validUrl(record.Url) {
			demandMap[record.DemandId] = record.Url
			beego.Debug(demandMap[record.DemandId])
		}
	}

	return demandMap, nil
}

//todo
func validUrl(url string) bool {
	return true
}
