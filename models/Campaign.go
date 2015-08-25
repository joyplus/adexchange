package models

//import (
//	"errors"
//	"fmt"
//	"reflect"
//	"strings"
//	"time"

//	"github.com/astaxie/beego/orm"
//)

type PmpCampaignCreative struct {
	Id          int    `orm:"column(id);auto"`
	Name        string `orm:"column(name);size(45);null"`
	StartDate   string `orm:"column(start_date);type(date);null"`
	EndDate     string `orm:"column(end_date);type(date);null"`
	Width       int    `orm:"column(width);null"`
	Height      int    `orm:"column(height);null"`
	CreativeUrl string `orm:"column(creative_url);size(255);null"`
	Status      int    `orm:"column(status);null"`
	LandingUrl  string `orm:"column(landing_url);size(500);null"`
}
