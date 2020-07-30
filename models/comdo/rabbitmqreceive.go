package comdo

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

//广播   //定时请求 或者常连接
func RabbitMqReceiveVn() {
	//接收参数
	Exchange := "vhost_test" // 交换机
	ExchangeType := "fanout" //交换机类型
	Queue := "test_logs"     //队列名

	// 建立链接
	conn, err := amqp.Dial("amqp://admin:123456@139.224.117.139:5672/vhost_test")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 声明一个主要使用的 exchange
	err = ch.ExchangeDeclare(
		Exchange,     // name
		ExchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// 声明一个常规的队列, 其实这个也没必要声明,因为 exchange 会默认绑定一个队列
	q, err := ch.QueueDeclare(
		Queue, // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,   // queue name, 这里指的是 test_logs
		"",       // routing key
		Exchange, // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	// 这里监听的是 test_logs
	msgs, err := ch.Consume(
		Queue, // queue name, 这里指的是 test_logs
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
			fmt.Println(string(d.Body))
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

//一对一
func RabbitMqReceiveV1() {
	// 建立链接
	conn, err := amqp.Dial("amqp://" + RabbitMQ_Test + "/" + ChatVirtualHost)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 这里监听的是 test_logs
	msgs, err := ch.Consume(
		MqQueueName, // queue name, 这里指的是 test_logs
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			//fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
			fmt.Println(string(d.Body))
		}
	}()

	//log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
