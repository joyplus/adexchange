package models

import (
	"adexchange/tools"
	"github.com/astaxie/beego"
)

type AvbDemand struct {
	AllocationId        int
	PmpAdspaceId        int
	DemandAdspaceId     int
	PmpAdspaceKey       string
	DemandAdspaceKey    string
	PlanImp             int
	PlanClk             int
	ActualImp           int
	ActualClk           int
	detailAllocationFlg bool
	allocationDetailMap map[string]bool
}

type AllocationDetail struct {
	Id            int
	AllocationId  int
	TargetingType string
	TargetingCode string
	PlanImp       int
	PlanClk       int
	ActualImp     int
	ActualClk     int
}

//func (this *AvbDemand) GenerateCommonResponse() CommonResponse {
//	res := CommonResponse{}
//	res.StatusCode = this.StatusCode
//	res.AdspaceKey = this.AdspaceKey
//	res.Bid = this.Bid
//	res.Adunit = this.Adunit

//	return res
//}

func (this *AvbDemand) CheckAvailable(adRequest *AdRequest) (avbFlg bool, targetingCode string) {
	beego.Debug(this)

	if this.detailAllocationFlg {
		provinceCode, cityCode := tools.QueryIP(adRequest.Ip)
		beego.Debug("Get location:" + provinceCode + cityCode)
		allocationKey1 := "PROVINCE" + "_" + provinceCode
		allocationKey2 := "CITY" + "_" + cityCode

		_, ok1 := this.allocationDetailMap[allocationKey1]

		if ok1 {
			avbFlg = true
			targetingCode = provinceCode
			return
		}

		_, ok2 := this.allocationDetailMap[allocationKey2]
		if ok2 {
			avbFlg = true
			targetingCode = cityCode
			return
		}

		//if ok1 || ok2 {
		//	avbFlg = true
		//}

	} else {
		if this.PlanImp > this.ActualImp {
			avbFlg = true
		}
	}
	return
}

func (this *AvbDemand) SetDetailAllocation(detail *AllocationDetail) {

	if this.allocationDetailMap == nil {
		this.detailAllocationFlg = true
		this.allocationDetailMap = make(map[string]bool)
	}
	if detail.PlanImp > detail.ActualImp {
		allocationKey := detail.TargetingType + "_" + detail.TargetingCode
		this.allocationDetailMap[allocationKey] = true
	}

	//this.allocationDetailMap["DETAIL_PLAN_IMP_"+midKey] = detail.PlanImp
	//this.allocationDetailMap["DETAIL_PLAN_CLK_"+midKey] = detail.PlanClk
	//this.allocationDetailMap["DETAIL_ACTUAL_IMP_"+midKey] = detail.AcutalImp
	//this.allocationDetailMap["DETAIL_ACTUAL_CLK_"+midKey] = detail.ActualClk

}

//record.PlanImp > record.ActualImp
