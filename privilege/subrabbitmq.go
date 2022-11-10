package privilege

import (
	"errors"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type rabbitmqManage struct {
	manage
	addr         string
	conn         *amqp.Connection
	channel      *amqp.Channel
	queueName    string // 队列名
	routingKey   string // key名
	exchangeName string // 交换机名
	exchangeType string // 交换机类型
}

func (m *rabbitmqManage) publish(key, data string) (bool, error) {
	m.queueName = key

	if m.channel == nil {
		m.connect(m.addr)
	}

	m.declareExchange()

	m.declareQueue()

	err := m.channel.Publish(m.exchangeName, "", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(data),
	})

	if err != nil {
		return false, errors.New("发布消息失败")
	}

	log.Println("发布成功")
	return true, nil
}

func (m *rabbitmqManage) sublish(key string, recevier chan string) {
	m.queueName = key

	if m.channel == nil {
		m.connect(m.addr)
	}

	m.declareExchange()

	m.declareQueue()

	err := m.channel.QueueBind(m.queueName, "", m.exchangeName, false, nil)
	if err != nil {
		panic(err)
	}

	messages, err := m.channel.Consume(m.queueName, "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	b := make(chan bool)
	go func() {
		for msg := range messages {
			recevier <- string(msg.Body)
		}
	}()
	<-b
}

func newRabbitmqManage(c ManageConfig) *rabbitmqManage {
	m := new(rabbitmqManage)
	m.addr = fmt.Sprintf("amqp://%s:%s@%s:%s", c.UserName, c.Password, c.Addr, c.Port)
	m.exchangeName = "test_exchange"
	m.exchangeType = "fanout"
	m.connect(m.addr)

	return m
}

// 连接
func (m *rabbitmqManage) connect(addr string) {
	// 创建连接
	conn, err := amqp.Dial(addr)
	if err != nil {
		log.Println("rabbimq创建连接失败, ", err)
		panic(err)
	}

	m.conn = conn

	// 打开通道
	if m.channel, err = conn.Channel(); err != nil {
		log.Println("rabbimq打开通道失败, ", err)
		panic(err)
	}
}

// 声明队列
func (m *rabbitmqManage) declareQueue() {
	if _, err := m.channel.QueueDeclare(m.queueName, false, false, false, false, nil); err != nil {
		log.Println("声明队列失败, ", err)
		panic(err)
	}
}

// 声明交换机
func (m *rabbitmqManage) declareExchange() {
	err := m.channel.ExchangeDeclare(m.exchangeName, m.exchangeType, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
}

// 关闭请求
func (m *rabbitmqManage) close() {
	if m.channel != nil {
		if err := m.channel.Close(); err != nil {
			log.Println("rabbimq关闭通道失败, ", err)
			panic(err)
		}
	}

	if m.conn != nil {
		if err := m.conn.Close(); err != nil {
			log.Println("rabbimq关闭链接失败, ", err)
			panic(err)
		}
	}
}
