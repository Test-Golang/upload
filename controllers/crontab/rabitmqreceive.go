package crontab

import (
	"fmt"
	"rabbitmq/models/comdo"

	"github.com/astaxie/beego"
)

var (
	ChatVirtualHost string
	MqQueueName     string
	RabbitMQ_Test   string
)

func Run() {
	if MqQueueName != "" && RabbitMQ_Test != "" && ChatVirtualHost != "" {
		comdo.RabbitMqReceiveV1()
	}
}
func init() {
	ChatVirtualHost = beego.AppConfig.String("RabbitMQ::ChatVirtualHost")
	if ChatVirtualHost == "" {
		comdo.LogError("ChatVirtualHost未配置")
		fmt.Println("ChatVirtualHost未配置")
	}
	MqQueueName = beego.AppConfig.String("RabbitMQ::MqQueueName")
	if MqQueueName == "" {
		comdo.LogError("MqQueueName未配置")
		fmt.Println("MqQueueName未配置")
	}
	RabbitMQ_Test = beego.AppConfig.String("RabbitMQ::RabbitMQTest")
	if RabbitMQ_Test == "" {
		comdo.LogError("RabbitMQ_Test未配置")
		fmt.Println("RabbitMQ_Test未配置")
	}
}
