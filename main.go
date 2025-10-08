package main

import (
	"myapp/common/middlewares"
	"myapp/config"
	routes_item "myapp/module/item/routes"
	routes_order "myapp/module/order/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectMongo()

	r := gin.Default()
	r.Use(middlewares.TrimJSONMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	routes_item.ItemRoutes(r)
	routes_order.OrderRoutes(r)
	r.Run(":8080")
}
