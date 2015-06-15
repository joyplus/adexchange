package models

const MAX_DEMAND_INDEX = 5

type PmpDemandInfo struct {
	AdspaceKey string
	DemandIds  []int
	index      int
}

func (this *PmpDemandInfo) AddDemand(demandId int) {
	if this.index == MAX_DEMAND_INDEX {
		panic("Reach the max demand index")
	}
	this.DemandIds[this.index] = demandId
	this.index++
}

func (this *PmpDemandInfo) InitDemand() {
	this.DemandIds = make([]int, MAX_DEMAND_INDEX)
	this.index = 0
}

func (this *PmpDemandInfo) GetDemandIds() []int {
	return this.DemandIds[:this.index]
}
