package routers

import (
	"rabbitmq/controllers/redis_test"

	"github.com/astaxie/beego"
)

func init() {
	//beego.Router("/mq/SendRabbitmq", &rabbitmq.Send{}, "post:SendRabbitmq")
	//修改推送状态
	beego.Router("/mq/UpMqState", &redis_test.RedisTest{}, "post:UpMqState")
}
