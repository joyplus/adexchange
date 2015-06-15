package engine

import (
	m "adexchange/models"
	"time"
)

func ScheduleInit(int minutes) {
	timer := time.NewTicker(time.Minute * minutes)
	for {
		select {
		case <-timer.C:
			m.InitEngineData()
		}
	}
}
