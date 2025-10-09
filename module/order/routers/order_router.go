package routers

import (
	"myapp/module/order/handlers"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(r *gin.Engine) {
	r.GET("/orders", handlers.GetAllOrders)
	r.GET("/orders/user/:userID", handlers.GetAllOrdersByCustomer)
	r.POST("/orders", handlers.CreateOrder)
	r.GET("/orders/:id", handlers.GetOrderByID)
	r.PATCH("/orders/info/:id", handlers.UpdateInfoOrder)
	r.PATCH("/orders/status/:id/:status", handlers.UpdatePendingOrder)
	r.PATCH("/orders/completed/:id", handlers.UpdateConfirmOrder)
}
