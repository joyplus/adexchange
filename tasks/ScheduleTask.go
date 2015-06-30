package tasks

import (
	"adexchange/engine"
	m "adexchange/models"
	"github.com/astaxie/beego"
	"time"
)

func InitEngineData() {

	beego.Debug("Start Init Engine Data")
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

	engine.SetupAdspaceMap(adspaceMap)
	engine.SetupAdspaceDemandMap(adspaceDemandMap)
	engine.SetupDemandMap(demandMap)

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

func GetAvbDemand() {

	beego.Debug("Check avalabile demand")
	avbDemandMap, err := m.GetAvbDemandMap(time.Now().Format("2006-01-02"))

	if err != nil {
		beego.Error(err.Error())
		return
	}

	engine.SetupAvbAdspaceDemandMap(avbDemandMap)

}
