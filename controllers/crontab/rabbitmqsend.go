package crontab

import (
	"fmt"
	"rabbitmq/models/comdo"
	"rabbitmq/models/detection"
	"time"

	"github.com/garyburd/redigo/redis"
)

func CrontabRun() {
	fmt.Println("定时任务启动...")
	c := comdo.NewWithSeconds()
	c.AddFunc("0/1 * * * * ?", CrontabSend) //每秒一次
	c.Start()
}

var Count int

func CrontabSend() {
	// h := time.Now().Hour()
	m := time.Now().Minute()
	s := time.Now().Second()
	openTime := time.Now().Format("2006-01-02 15:04:05")
	rc := comdo.GetRedisPool("test")
	defer rc.Close()
	rc.Do("SELECT", 0)
	state, _ := redis.String(rc.Do("GET", "test"))
	if state != "1" {
		return
	}
	switch s {
	case 0: //每分钟一次
		SendRabbitmq1("我还活着,没死...")
		if m == 0 { //一小时一次
			SendRabbitmq1("正在插入数据啦...")
			//插入数据库记录
			Number := detection.SellDetection()
			if Number.Num == 0 {
				detection.InlDetection()
			} else {
				if Count == 0 {
					detection.InlDetections(Number.Num)
				} else {
					detection.UplDetection(Count, Number.MaxId)
				}
			}
			SendRabbitmq1("数据已经插入啦...")
			if Count == 0 {
				Count += 1
			}
		}
	default: //一秒一次
		SendRabbitmq1("现在时间：" + openTime)
	}
}

//推送 一对一
func SendRabbitmq1(parameter string) {
	comdo.InitRabbitMQPerpetual(parameter) //(队列名   内容)
}

//推送 广播
func SendRabbitmq2(parameter string) {
	comdo.InitRabbitMQBroadcast(comdo.ChatVirtualHost, parameter) //交换机名  内容
}

//延时推送
func SendRabbitmq3(parameter string) {
	comdo.RabbitMQToDLX("letter", parameter, 100) //(队列名   内容)
}
