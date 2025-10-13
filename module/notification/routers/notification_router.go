package routers

import (
	"myapp/module/notification/handler"

	"github.com/gin-gonic/gin"
)

func NotificationRoutes(r *gin.Engine) {
	r.GET("/ws/orders", handler.OrdersWSHandler)
}
