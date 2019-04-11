package routers

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"rabbitmq/controllers"
)


func SetRouter(router iris.Party, path string) iris.Party {
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowedMethods:   []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
	})
	router = router.Party("/", crs).AllowMethods(iris.MethodOptions)
	r := router.Party(path)

	r.Get("/mq",controllers.Mq)//查询商品列表

	return router
}