package models

import (
	"github.com/astaxie/beego"
)

const MAX_DEMAND_INDEX = 10

type PmpDemandInfo struct {
	AdspaceKey string
	//DemandIds        []int
	AryDemandAdspaceKey []string
	index               int
	demandKeyMap        map[string]bool
}

//func (this *PmpDemandInfo) AddDemand(demandId int) {
//	if this.index == MAX_DEMAND_INDEX {
//		panic("Reach the max demand index")
//	}
//	this.DemandIds[this.index] = demandId
//	this.index++
//}

//func (this *PmpDemandInfo) InitDemand() {
//	this.DemandIds = make([]int, MAX_DEMAND_INDEX)
//	this.index = 0
//}

//func (this *PmpDemandInfo) GetDemandIds() []int {
//	return this.DemandIds[:this.index]
//}

func (this *PmpDemandInfo) AddDemandAdspace(demandAdspace string) {
	if this.index == MAX_DEMAND_INDEX {
		beego.Critical("Reach the max demand index")
	}
	_, ok := this.demandKeyMap[demandAdspace]
	if !ok {
		this.AryDemandAdspaceKey[this.index] = demandAdspace
		this.demandKeyMap[demandAdspace] = true
		this.index++
	}

}

func (this *PmpDemandInfo) InitDemandAdspace() {
	this.AryDemandAdspaceKey = make([]string, MAX_DEMAND_INDEX)
	this.demandKeyMap = make(map[string]bool)
	this.index = 0
}

func (this *PmpDemandInfo) GetDemandAdspaceKeys() []string {
	return this.AryDemandAdspaceKey[:this.index]
}
