package main

import (
	"myapp/common/middlewares"
	"myapp/config/db"
	"myapp/config/oauth"
	item_handler "myapp/module/item/handlers"
	notification_handler "myapp/module/notification/handlers"
	order_handler "myapp/module/order/handlers"
	user_handler "myapp/module/user/handlers"
	"os"
	"strings"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectMongo()
	oauth.InitGoogle()

	r := gin.Default()
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
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
	ItemRoutes(r, authGroup)
	OrderRoutes(r, authGroup)
	AuthRoutes(r, authGroup)
	FCMRouter(r, authGroup)

	r.Run(":8080")
}

func ItemRoutes(r *gin.Engine, a *gin.RouterGroup) {
	r.GET("/items", item_handler.GetAllItems)
	r.GET("/items/:id", item_handler.GetItemByID)

	a.POST("/items", item_handler.CreateItem)
	a.PUT("/items/:id", item_handler.UpdateItem)
	a.DELETE("/items/:id", item_handler.DeleteItem)
}

func OrderRoutes(r *gin.Engine, a *gin.RouterGroup) {
	r.GET("/orders", order_handler.GetAllOrders)
	r.POST("/orders", order_handler.CreateOrder)

	a.GET("/orders/user/:userID", order_handler.GetAllOrdersByCustomer)
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
