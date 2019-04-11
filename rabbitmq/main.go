package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	_ "rabbitmq/init"
	"rabbitmq/routers"
)

func main()  {
	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())
	routers.SetRouter(app, "/api/v1")
	app.Run(iris.Addr(":4010"), iris.WithoutPathCorrection)
}