package main

import (
	"myapp/common/middlewares"
	"myapp/config/db"
	"myapp/config/oauth"
	item_handler "myapp/module/item/handlers"
	notification_handler "myapp/module/notification/handlers"
	order_handler "myapp/module/order/handlers"
	user_handler "myapp/module/user/handlers"
	"myapp/module/user/models"
	"os"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectMongo()
	oauth.InitGoogle()

	r := gin.Default()

	corsConfig := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(corsConfig))

	//Middlewares
	r.Use(middlewares.TrimJSONMiddleware())
	adminGroup := r.Group("")
	adminGroup.Use(middlewares.AuthMiddleware())
	adminGroup.Use(middlewares.RequireRoles(string(models.ADMIN)))

	customerGroup := r.Group("")
	customerGroup.Use(middlewares.AuthMiddleware())
	customerGroup.Use(middlewares.RequireRoles(string(models.CUSTOMER)))

	//Middlewares-Auth
	adminGroup.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	//Router rest
	ItemRoutes(r, adminGroup)
	OrderRoutes(r, adminGroup, customerGroup)
	AuthRoutes(r, adminGroup)
	FCMRouter(r, adminGroup)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)

}

func ItemRoutes(r *gin.Engine, a *gin.RouterGroup) {
	r.GET("/items", item_handler.GetAllItems)
	r.GET("/items/:id", item_handler.GetItemByID)

	a.POST("/items", item_handler.CreateItem)
	a.PUT("/items/:id", item_handler.UpdateItem)
	a.DELETE("/items/:id", item_handler.DeleteItem)
}

func OrderRoutes(r *gin.Engine, a *gin.RouterGroup, c *gin.RouterGroup) {
	r.GET("/orders", order_handler.GetAllOrders)
	r.POST("/orders", order_handler.CreateOrder)

	c.GET("/orders/user/:userID", order_handler.GetAllOrdersByCustomer)
	
	a.GET("/orders/:id", order_handler.GetOrderByID)
	a.PATCH("/orders/info/:id", order_handler.UpdateInfoOrder)
	a.PATCH("/orders/status/:id/:status", order_handler.UpdatePendingOrder)
	a.PATCH("/orders/completed/:id", order_handler.UpdateConfirmOrder)
}

func AuthRoutes(r *gin.Engine, a *gin.RouterGroup) {
	r.GET("/auth/google", user_handler.GoogleLogin)
	r.GET("/auth/google/callback", user_handler.GoogleCallback)
	r.GET("/auth/introspect/:token", user_handler.InspectToken)
	r.POST("/auth/admin", user_handler.GenerateTokenAdmin)
}

func FCMRouter(r *gin.Engine, a *gin.RouterGroup) {
	r.POST("/test/send", notification_handler.SendNotification)
}
