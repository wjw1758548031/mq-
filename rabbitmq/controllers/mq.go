package controllers

import (
	"fmt"
	"github.com/kataras/iris"
	in "rabbitmq/init"
)

//查询商品列表
func Mq(ctx iris.Context) {

	fmt.Println("进入mq接口")
	newMQ := in.RabbitMQ{}
	//线连接到具体的mq
	mq := newMQ.Newmq()
	//在发送消息 这条消息是个连接有所pay的交换机发送的
	err := mq.Send("王建文")
	fmt.Println("err:",err)

}
