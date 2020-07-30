package detection

import (
	"rabbitmq/models/comdo"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

type Detection struct {
	Id        int    `json:"id" form:"id"`
	BeginTime string `json:"begin_time" form:"begin_time"`
	DownTime  string `json:"down_time" form:"down_time"`
	HzMinute  int    `json:"hz_minute" form:"hz_minute"`
	TheNum    int    `json:"the_num" form:"the_num"`
}

type Statistics struct {
	Num   int `json:"num" form:"num"`
	MaxId int `json:"max_id" form:"max_id"`
}

func SellDetection() (Statistics Statistics) {
	o := orm.NewOrm()
	o.Using("detection")
	sql := "SELECT count(*) as num,max(id) as max_id FROM `l_detection`"
	err := o.Raw(sql).QueryRow(&Statistics)
	if err != nil {
		comdo.LogError("SellDetection: %s %s", err.Error(), sql)
	}
	return Statistics
}

func InlDetection() {
	o := orm.NewOrm()
	o.Using("detection")
	BeginTime := time.Now().Format("2006-01-02 15:04:05")
	sql := "INSERT INTO `l_detection` (`begin_time`, `down_time`, `hz_minute`, `the_num`) VALUES ( '" + BeginTime + "', 0,1,1);"
	_, err := o.Raw(sql).Exec()
	if err != nil {
		comdo.LogError("InlDetection: %s %s", err.Error(), sql)
	}
}

func InlDetections(Number int) {
	o := orm.NewOrm()
	o.Using("detection")
	BeginTime := time.Now().Format("2006-01-02 15:04:05")
	sql := "INSERT INTO `l_detection` (`begin_time`, `down_time`, `hz_minute`, `the_num`) VALUES ( '" + BeginTime + "', " + strconv.Itoa(Number) + ",1,1);"
	_, err := o.Raw(sql).Exec()
	if err != nil {
		comdo.LogError("InlDetection: %s %s", err.Error(), sql)
	}
}

func UplDetection(Count, MaxId int) {
	o := orm.NewOrm()
	o.Using("detection")
	var up string
	sql := "UPDATE `l_detection` SET the_num = (the_num + 1 )" + up + " WHERE id = " + strconv.Itoa(MaxId)
	if Count == 0 {
		up += ",begin_time = " + time.Now().Format("2006-01-02 15:04:05")
		up += ",down_time = (down_time + 1 )"
	}
	_, err := o.Raw(sql).Exec()
	if err != nil {
		comdo.LogError("InlDetection: %s %s", err.Error(), sql)
	}
}
