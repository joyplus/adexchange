package engine

import (
	m "adexchange/models"
	"github.com/astaxie/beego"
)

func mockupGeoLocation(adRequest *m.AdRequest, targetingCode string) {
	beego.Debug("Start to setup geo location.")
}
