package websocket

import (
	"fmt"
	"net/url"
	"rabbitmq/models/comdo"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

//定义连接的服务端的网址
//var addr = flag.String("addr", "139.224.117.139:8000", "http service address")
var WebsocketTest string

func WebsocketTestApi() {
	// u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	u := url.URL{Scheme: "ws", Host: WebsocketTest, Path: "/ws"}
	var dialer *websocket.Dialer

	//通过Dialer连接websocket服务器
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//go timeWriter(conn)
	//打印接收到的消息或者错误

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		fmt.Printf("received: %s\n", message)
	}
}

//func timeWriter(conn *websocket.Conn) {
//	for {
//		time.Sleep(time.Second * 2)
//		conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
//	}
//}

func init() {
	WebsocketTest = beego.AppConfig.String("websocket::WebsocketTest")
	if WebsocketTest == "" {
		comdo.LogError("WebsocketTest未配置")
		fmt.Println("WebsocketTest未配置")
	}
}
