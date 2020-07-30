package comdo

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/robfig/cron"
)

var Log *logs.BeeLogger = nil
var btime string

//错误日志
func LogError(format string, v ...interface{}) {
	Log.EnableFuncCallDepth(true)
	Log.Error(format, v...)
}

//日志处理
func init() {
	Log = logs.NewLogger(10000)
	if beego.AppConfig.String("runmode") == "dev" {
		Log.SetLogger("console", "")
	}

	appname := beego.AppConfig.String("appname")

	Log.Async(2000)
	Log.SetLogger("file", `{"filename":"logs/`+appname+`.log","maxlines":0,"maxsize":0,"daily":true,"maxdays":7}`)

	Log.SetLogFuncCallDepth(3)
	Log.EnableFuncCallDepth(true)

	btime = beego.AppConfig.String("Estime")
}

//返回一个支持至 秒 级别的 cron
func NewWithSeconds() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}
