package comdo

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"github.com/garyburda/redigoa/redisa"
)

var (
	REDIS_POOL_TEST *redisa.Pool
	REDIS_TEST_HOST string
	REDIS_PWD       string
)

func init() {
	//测试
	REDIS_TEST_HOST = beego.AppConfig.String("redis::redis_test_host")
	REDIS_PWD = beego.AppConfig.String("redis::redis_pwd")

	if len(REDIS_TEST_HOST) > 1 {
		InitRedisPoolByTest()
	} else {
		LogError("EG_SERVER_REDIS_HOST 配置为空")
	}
}

//初始化连接池 REDIS_TEST_HOST
func InitRedisPoolByTest() {
	// 建立连接池
	REDIS_POOL_TEST = &redisa.Pool{
		MaxIdle:     1000,
		MaxActive:   5000,
		IdleTimeout: 30 * time.Second,
		Dial: func() (redisa.Conn, error) {
			c, err := redisa.Dial("tcp", REDIS_TEST_HOST)
			if err != nil {
				return nil, err
			}
			if len(REDIS_PWD) > 0 {
				if _, err := c.Do("AUTH", REDIS_PWD); err != nil {
					c.Close()
					LogError("REDIS AUTH %s", err.Error())
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redisa.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

//获取一个redis连接
func GetRedisPool(rtype string) redis.Conn {
	if rtype == "test" {
		return REDIS_POOL_TEST.Get()
	}
	return REDIS_POOL_TEST.Get()
}
