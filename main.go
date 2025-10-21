package main

import (
	"myapp/common/middlewares"
	"myapp/config/db"
	"myapp/config/oauth"
	routes_item "myapp/module/item/routes"
	routers_notification "myapp/module/notification/routers"
	routes_order "myapp/module/order/routers"
	routes_user "myapp/module/user/routers"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectMongo()
	oauth.InitGoogle()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//Middlewares
	r.Use(middlewares.TrimJSONMiddleware())
	authGroup := r.Group("")
	authGroup.Use(middlewares.AuthMiddleware())
	authGroup.Use(middlewares.RequireRoles("ADMIN"))
	//Middlewares-Auth
	authGroup.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	//Router rest
	routes_item.ItemRoutes(r)
	routes_order.OrderRoutes(r)
	routes_user.AuthRoutes(r)
	routers_notification.FCMRouter(r)

	r.Run(":8080")
}
