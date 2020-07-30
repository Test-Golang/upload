package main

import (
	"net/http"
	"rabbitmq/controllers/crontab"
	"rabbitmq/controllers/websocket"
	_ "rabbitmq/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	maxIdle := beego.AppConfig.DefaultInt("dbconfig::maxIdle", 20)
	maxConn := beego.AppConfig.DefaultInt("dbconfig::maxConn", 30)
	l_detection := beego.AppConfig.String("dbconfig::l_detection")
	orm.RegisterDataBase("default", "mysql", l_detection, maxIdle, maxConn)
	orm.RegisterDataBase("l_detection", "mysql", l_detection, maxIdle, maxConn)
}
func main() {
	beego.ErrorHandler("404", page_go)
	// orm.Debug = true //打开查询日志
	//长连接
	crontab.CrontabRun()
	//消费
	go crontab.Run()
	//-------------------------------------------------------------------分割线--------------------------------------------------------------------------
	websocket.WebsocketTestApi()
	beego.Run()
}
func page_go(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("搞错了,重来... 提示：404"))
}
