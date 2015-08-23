package tasks

import (
	"adexchange/engine"
	m "adexchange/models"
	"github.com/astaxie/beego"
	"time"
)

func InitEngineData() {

	beego.Info("Start Init Engine Data")
	adspaceMap, adspaceDemandMap, err := m.GetMatrixData()

	if err != nil {
		panic(err.Error())
	}
	beego.Debug(adspaceMap)
	beego.Debug(adspaceDemandMap)

	demandMap, err := m.GetDemandInfo()

	if err != nil {
		panic(err.Error())
	}

	beego.Debug(demandMap)

	pmpAdspaceMap, err := m.GetPmpInfo()

	if err != nil {
		panic(err.Error())
	}

	beego.Debug(pmpAdspaceMap)

	tplHashSet, _ := m.GetTplSet()

	engine.SetupAdspaceMap(adspaceMap)
	engine.SetupAdspaceDemandMap(adspaceDemandMap)
	engine.SetupDemandMap(demandMap)
	engine.SetupPmpAdspaceMap(pmpAdspaceMap)
	engine.SetTplHashSet(tplHashSet)

}

func ScheduleInit(minutes int) {

	timer := time.NewTicker(time.Minute * time.Duration(minutes))
	for {
		select {
		case <-timer.C:
			InitEngineData()
		}
	}
}

func CheckAvbDemandInit(minutes int) {

	timer := time.NewTicker(time.Minute * time.Duration(minutes))
	for {
		select {
		case <-timer.C:
			CheckAvbDemand()
		}
	}
}

func CheckAvbDemand() {

	beego.Debug("Check avalabile demand")
	avbDemandMap, err := m.GetAvbDemandMap(time.Now().Format("2006-01-02"))

	if err != nil {
		beego.Error(err.Error())
		return
	}

	beego.Debug(avbDemandMap)
	engine.SetupAvbAdspaceDemandMap(avbDemandMap)

}
