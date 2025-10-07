package main

import (
	"myapp/config"
	"myapp/module/item/routes"
	"myapp/module/item/services"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectMongo()
	services.InitCollections()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	routes.ItemRoutes(r)
	r.Run(":8080")
}
