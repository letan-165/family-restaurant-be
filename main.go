package main

import (
	"myapp/common/middlewares"
	"myapp/config/db"
	"myapp/config/oauth"
	routes_item "myapp/module/item/routes"
	routes_order "myapp/module/order/routers"
	routes_user "myapp/module/user/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectMongo()
	oauth.InitGoogle()

	r := gin.Default()
	r.Use(middlewares.TrimJSONMiddleware())

	authGroup := r.Group("")
	authGroup.Use(middlewares.AuthMiddleware())

	authGroup.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	routes_item.ItemRoutes(r)
	routes_order.OrderRoutes(r)
	routes_user.AuthRoutes(r)
	r.Run(":8080")
}
