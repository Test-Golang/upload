package comdo

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/streadway/amqp"
)

var (
	ChannelTest      *amqp.Channel //通道
	RabbitMQ_Test    string        //连接地址以及账号密码
	ChatVirtualHost  string        //交换机
	ChatExchangeName string
	ChatExchangeType string
	ChatRoutingKey   string
	MqQueueName      string
	MqTestType       = true
)

//一对一 一对多 队列  永久
func InitRabbitMQPerpetual(data string) {
	//检测连接状态 重连
	Reconnection()
	// fmt.Println("一对一已启动...队列名：s%", QueueName)
	udata, err := json.MarshalIndent(data, "", "  ") //加标签
	if err != nil {
		LogError("JSON失败:%s", err.Error())
		return
	}
	err = ChannelTest.Publish(
		ChatExchangeName, // exchange
		ChatRoutingKey,   // routing key
		false,            // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         udata,
		})
	if err != nil {
		LogError("Publish:%s", err.Error())
		return
	}
	return
}

//广播
func InitRabbitMQBroadcast(ExchangeName, data string) {
	//检测连接状态 重连
	Reconnection()
	//fmt.Println("广播已启动...交换机：s%", ExchangeName)
	udata, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		LogError("JSON失败:%s", err.Error())
		return
	}
	//defer ChannelTest.Close() //推迟关闭
	//声明交换机
	err = ChannelTest.ExchangeDeclare(
		ExchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		LogError("channel:%s", err.Error())
		return
	}
	err = ChannelTest.Publish(
		ExchangeName, // exchange
		"",           // routing key
		false,        // mandatory
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        udata,
		})

	if err != nil {
		LogError("Publish:%s", err.Error())
	} else {
		//LogError("MQ广播数据成功")
	}
}

// 死信队列（做延迟推送）
func RabbitMQToDLX(Queue, parameter string, Otime int32) error {
	select {
	case <-ChannelTest.NotifyClose(make(chan *amqp.Error)):
		//LogError("rabbitmq is closed")
		InitRabbitMQChatConnect()
	default:
		//LogError("rabbitmq is connected")
	}
	udata, err := json.MarshalIndent(parameter, "", "  ")
	if err != nil {
		LogError("JSON失败:%s", err.Error())
		return err
	}
	// fmt.Println("延迟（死信）推送已启动...队列名：s%", Queue)
	_, err = ChannelTest.QueueDeclare(
		Queue,
		false,
		false,
		false,
		false,
		amqp.Table{
			// 当消息过期时把消息发送到 logs 这个 exchange
			"x-dead-letter-exchange":    ChatExchangeName,
			"x-dead-letter-routing-key": ChatRoutingKey,
			"x-expires":                 Otime + 3000, //超过这个时间 发送到上面交换机
		})
	err = ChannelTest.Publish(
		"",    // exchange
		Queue, // routing key
		false, // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         udata,
			Expiration:   fmt.Sprintf("%d", Otime), //过期时间
		})
	if err != nil {
		LogError("Publish:%s", err.Error())
		return err
	}
	return nil
}

//MQ检测状态 重连
func Reconnection() {
	//方案1
	if ChannelTest.IsClosed() {
		//log.Println("连接断开，重新连接..f.")
		InitRabbitMQChatConnect()
	}

	//方案2
	// select {
	// case <-ChannelTest.NotifyClose(make(chan *amqp.Error)):
	// 	//LogError("rabbitmq is closed")
	// 	InitRabbitMQChatConnect()
	// default:
	// 	//LogError("rabbitmq is connected")
	// }
}

//连接MQ
func InitRabbitMQChatConnect() {
	if !MqTestType {
		LogError("配置错误:Mq启动失败...")
		return
	}
	connection, err := amqp.Dial("amqp://" + RabbitMQ_Test + "/" + ChatVirtualHost)
	if err != nil {
		LogError("Dial:%s %s", "amqp://"+RabbitMQ_Test+"/"+ChatVirtualHost, err.Error())
		return
	}
	//连接通道
	ChannelTest, err = connection.Channel()
	if err != nil {
		LogError("channel:%s", err.Error())
		return
	}
	//申明交换机
	err = ChannelTest.ExchangeDeclare(
		ChatExchangeName, // name
		ChatExchangeType, // type   //交换机类型
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		LogError("channel:%s", err.Error())
		return
	}

	//申明队列
	q, err := ChannelTest.QueueDeclare(
		MqQueueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		LogError("channel:%s", err.Error())
		return
	}

	// 队列和交换机绑定
	err = ChannelTest.QueueBind(
		q.Name,         // queue name  队列的名字
		ChatRoutingKey, // routing key  广播模式不需要这个
		// "logs", // exchange  交换机名字
		ChatExchangeName, // exchange  交换机名字
		false,
		nil)
	if err != nil {
		LogError("channel:%s", err.Error())
		return
	}
}

func init() {
	ChatVirtualHost = beego.AppConfig.String("RabbitMQ::ChatVirtualHost")
	if ChatVirtualHost == "" {
		MqTestType = false
		LogError("ChatVirtualHost未配置")
		fmt.Println("ChatVirtualHost未配置")
	}

	ChatRoutingKey = beego.AppConfig.String("RabbitMQ::ChatRoutingKey")
	if ChatRoutingKey == "" {
		MqTestType = false
		LogError("ChatRoutingKey未配置")
		fmt.Println("ChatRoutingKey未配置")
	}

	ChatExchangeType = beego.AppConfig.String("RabbitMQ::ChatExchangeType")
	if ChatExchangeType == "" {
		MqTestType = false
		LogError("ChatExchangeType未配置")
		fmt.Println("ChatExchangeType未配置")
	}

	ChatExchangeName = beego.AppConfig.String("RabbitMQ::ChatExchangeName")
	if ChatExchangeName == "" {
		MqTestType = false
		LogError("ChatExchangeName未配置")
		fmt.Println("ChatExchangeName未配置")
	}

	MqQueueName = beego.AppConfig.String("RabbitMQ::MqQueueName")
	if MqQueueName == "" {
		MqTestType = false
		LogError("MqQueueName未配置")
		fmt.Println("MqQueueName未配置")
	}

	RabbitMQ_Test = beego.AppConfig.String("RabbitMQ::RabbitMQTest")
	if RabbitMQ_Test == "" {
		MqTestType = false
		LogError("RabbitMQ_Test未配置")
		fmt.Println("RabbitMQ_Test未配置")
	}
	if MqTestType {
		InitRabbitMQChatConnect()
	}
}
