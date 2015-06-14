package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type PmpAdplace struct {
	Id          int       `orm:"column(id);auto"`
	Name        string    `orm:"column(name);size(255)"`
	Description string    `orm:"column(description);size(500);null"`
	DelFlg      int8      `orm:"column(del_flg);null"`
	CreateUser  int       `orm:"column(create_user);null"`
	CreateTime  time.Time `orm:"column(create_time);type(timestamp);null"`
	UpdateUser  int       `orm:"column(update_user);null"`
	UpdateTime  time.Time `orm:"column(update_time);type(timestamp);null"`
	AdplaceKey  string    `orm:"column(adplace_key);size(50);null"`
	SecretKey   string    `orm:"column(secret_key);size(50);null"`
	EstDailyImp int       `orm:"column(est_daily_imp);null"`
	EstDailyClk int       `orm:"column(est_daily_clk);null"`
	EstDailyCtr float32   `orm:"column(est_daily_ctr);null"`
}

func init() {
	orm.RegisterModel(new(PmpAdplace))
}

// AddPmpAdplace insert a new PmpAdplace into database and returns
// last inserted Id on success.
func AddPmpAdplace(m *PmpAdplace) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetPmpAdplaceById retrieves PmpAdplace by Id. Returns error if
// Id doesn't exist
func GetPmpAdplaceById(id int) (v *PmpAdplace, err error) {
	o := orm.NewOrm()
	v = &PmpAdplace{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPmpAdplace retrieves all PmpAdplace matches certain condition. Returns empty list if
// no records exist
func GetAllPmpAdplace(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(PmpAdplace))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []PmpAdplace
	qs = qs.OrderBy(sortFields...)
	if _, err := qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdatePmpAdplace updates PmpAdplace by Id and returns error if
// the record to be updated doesn't exist
func UpdatePmpAdplaceById(m *PmpAdplace) (err error) {
	o := orm.NewOrm()
	v := PmpAdplace{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePmpAdplace deletes PmpAdplace by Id and returns error if
// the record to be deleted doesn't exist
func DeletePmpAdplace(id int) (err error) {
	o := orm.NewOrm()
	v := PmpAdplace{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&PmpAdplace{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
