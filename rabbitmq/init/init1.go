package init

import (
	"fmt"
)

func init() {
	newMQ := RabbitMQ{}
	mq := newMQ.Newmq()
	msg, err := mq.Subscribe("jieshouji1")
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		for {
			data, _ := <-msg
			if data.Body != nil {
				fmt.Println("收到消息", string(data.Body))
			}
		}
	}()

	init1()
}







func init1() {
	newMQ := RabbitMQ{}
	//创建连接和交换机
	mq := newMQ.Newmq()
	//绑定交换机并返回chan消息
	msg, err := mq.Subscribe("jieshouji2")
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		for {
			data, _ := <-msg
			if data.Body != nil {
				fmt.Println("收到消息", string(data.Body))
			}
		}
	}()
}