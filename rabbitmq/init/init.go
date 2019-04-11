package init

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"rubik/server/common"
	"time"
)

const rabbitMQ string = "amqp://localhost:5672/"

type RabbitMQ struct {
	Url string
	Exchange string
	RouterKey string
	Kind	string
	Conn *amqp.Connection
	Channel *amqp.Channel
	done          chan bool
	notifyClose   chan *amqp.Error
	publishArray	  commons.ArrayString
}

//var Mq RabbitMQ

func (this *RabbitMQ) Newmq() RabbitMQ{
		mq := RabbitMQ{Url:"amqp://localhost:5672/", Exchange:"pay", RouterKey:"", Kind:"fanout"}
		if err := mq.Connect(); err != nil {
			panic(err)
		}
		mq.notifyClose = make(chan *amqp.Error)
		return mq
}

func (this *RabbitMQ) Connect() (err error){
	if this.Conn, err = amqp.Dial(this.Url); err != nil {
		log.Println("连接RabbitMQ "+this.Url+"失败")
		return
	}
	if this.Channel, err = this.Conn.Channel(); err != nil{
		return
	}
	log.Println("连接RabbitMQ成功")
	err = this.CreateExchange(this.Exchange)
	return nil
}


//创建交换机
func (this *RabbitMQ) CreateExchange(exchange string) (err error) {
	if exchange == "" {
		return
	}
	if err = this.Channel.ExchangeDeclare(
		exchange, // name
		this.Kind, // type
		true,     // 持久
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	); err != nil {
		return
	}
	return
}


func (this *RabbitMQ) Subscribe(name string) (<-chan amqp.Delivery, error){
	//定义队列
	q, err := this.Channel.QueueDeclare(
		name,    // name
		true, // 持久
		false, // delete when unused
		false,  // 独占 关闭时自动删除队列
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}
	//将交换机和对列绑定
	if err := this.Channel.QueueBind(
		q.Name, // queue name
		this.RouterKey,     // routing key
		this.Exchange, // exchange
		false,
		nil); err != nil {
		return nil, err
	}

	msgs, err := this.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, err
	}
	return msgs, err
}


func (this *RabbitMQ) Send(content string) (err error) {
	if err := this.Channel.Publish(
		this.Exchange, // exchange
		this.RouterKey,     // routing key
		false,  // 强制
		false,  // 立即
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(content),
		}); err != nil {
		fmt.Println("PublishOne", err)
		//加入发送队列，防止断网
		//this.publishArray.Append(content)
		time.Sleep(time.Second*3)
		this.Connect()
		this.Send(content)
		return nil
	}
	log.Println("发布消息成功", this.RouterKey, content)
	return
}




//发布
func (this *RabbitMQ) LoopPublish() {
	go func() {
		for {
			publishArray := this.publishArray.Get()
			for k, v := range publishArray{
				err := this.Send(v)
				if err != nil {
					publishArray = append(publishArray, v)
					time.Sleep(time.Millisecond*100)
				}
				//log.Println("发布消息成功", this.RouterKey, v)
				this.publishArray.Slice(k)
				break
			}
			if len(publishArray) == 0 {
				time.Sleep(time.Millisecond*100)
			}
		}
	}()
}

func (this *RabbitMQ) AddPublish(content string) {
	//加入发送队列，防止断网
	this.publishArray.Append(content)
}

func (this *RabbitMQ) AddPublishs(contents []string) {
	//加入发送队列，防止断网
	for _, v := range contents {
		this.publishArray.Append(v)
	}
}

func (this *RabbitMQ) Publish(content string) error {
	err := this.Send(content)
	return err
}
