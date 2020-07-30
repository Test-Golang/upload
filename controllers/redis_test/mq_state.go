package redis_test

import (
	"fmt"
	"rabbitmq/controllers"
	"rabbitmq/models/comdo"

	"github.com/garyburd/redigo/redis"
)

type RedisTest struct {
	controllers.Base
}

//利用redis控制MQ 的状态
func (test *RedisTest) UpMqState() {
	rc := comdo.GetRedisPool("test")
	defer rc.Close()
	rc.Do("SELECT", 0)
	State, _ := redis.String(rc.Do("GET", "test"))
	var states string
	if State == "1" {
		State = "0"
		states = "推送已关闭"
	} else {
		State = "1"
		states = "推送已开启"
	}
	rc.Do("SET", "test", State)
	fmt.Println(states)
	test.EchoJSON("1101", states, "")
}
