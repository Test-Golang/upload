package controllers

import (
	"net/url"
	"rabbitmq/models/comdo"
	"regexp"
	"strings"

	"github.com/astaxie/beego"
)

type Base struct {
	beego.Controller
}

var ReStr = `(?:')|(?:")|(?:#)|(/\\*(?:.|[\\n\\r])*?\\)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`

func (this *Base) BgetString(data string) string {
	pdata := strings.TrimSpace(this.GetString(data))
	if len(pdata) > 0 {
		if this.ReString(pdata) == true {
			return ""
		}
	}
	return pdata
}

func (this *Base) ReString(pdata string) bool {
	pdata = strings.ToLower(pdata)

	pdatas, err := url.QueryUnescape(pdata)
	if err != nil {
		comdo.LogError("QueryUnescape error:  %s %s", err.Error(), pdata)
		return false
	}
	re, err := regexp.Compile(ReStr)
	if err != nil {
		comdo.LogError("ReString error:  %s", err.Error())
		return false
	}
	return re.MatchString(pdatas)
}

func (this *Base) EchoJSON(code, msg string, data interface{}) {
	jsonData := this.GetJSON(code, msg, data)
	this.Data["json"] = jsonData
	this.ServeJSON()
}
func (this *Base) GetJSON(code, msg string, data interface{}) map[string]interface{} {
	var jsonData = make(map[string]interface{})
	jsonData["Code"] = code
	jsonData["Msg"] = msg
	jsonData["Data"] = data
	return jsonData
}
