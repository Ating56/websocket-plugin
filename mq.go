package websocketplugin

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MQconf struct {

	// Url RabbitMQ连接地址
	Url string

	// Exchange RabbitMQ交换机名称
	Exchange string

	// Queue RabbitMQ队列名称
	Queue string
}

type RabbitMQ struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
	MQconf
}

var GlobalMQInstance *RabbitMQ

/*
 * InitMQ
 * 初始化RabbitMQ，建立连接，创建实例
 * @param conf RabbitMQ配置信息
 */
func InitMQ(conf MQconf) (*RabbitMQ, error) {
	conn, err := amqp.Dial(conf.Url)
	if err != nil {
		log.Println("InitMQ amqp.Dial err:", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Println("InitMQ conn.Channel err:", err)
		return nil, err
	}

	GlobalMQInstance = &RabbitMQ{
		Conn:   conn,
		Ch:     ch,
		MQconf: conf,
	}

	return GlobalMQInstance, nil
}

/*
 * Publish
 * 发布消息到RabbitMQ交换机
 * @param message 消息内容，[]byte格式
 */
func (rmq *RabbitMQ) Publish(message []byte) error {
	// ExchangeDeclare
	err := rmq.Ch.ExchangeDeclare(
		rmq.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Publish ExchangeDeclare err:", err)
		return err
	}

	// Publish
	err = rmq.Ch.Publish(
		rmq.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	if err != nil {
		log.Println("Publish Publish err:", err)
		return err
	}

	return nil
}

/*
 * Consume
 * 从RabbitMQ队列消费消息
 */
func (rmq *RabbitMQ) Consume() {
	// ExchangeDeclare
	err := rmq.Ch.ExchangeDeclare(
		rmq.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Consume ExchangeDeclare err:", err)
		return
	}

	// QueueDeclare
	q, err := rmq.Ch.QueueDeclare(
		rmq.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Consume QueueDeclare err:", err)
		return
	}

	// QueueBind
	err = rmq.Ch.QueueBind(
		q.Name,
		"",
		rmq.Exchange,
		false,
		nil,
	)
	if err != nil {
		log.Println("Consume QueueBind err:", err)
		return
	}

	// Consume
	message, err := rmq.Ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Consume Consume err:", err)
		return
	}

	forever := make(chan struct{})
	go func() {
		for i := range message {
			asyncStoreInMongo(i.Body)
		}
	}()
	<-forever
}
