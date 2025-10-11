package routers

import (
	"myapp/module/notification/ws"

	"github.com/gin-gonic/gin"
)

func NotificationRoutes(r *gin.Engine) {
	r.GET("/ws/orders", func(c *gin.Context) {
		ws.OrdersWSHandler(c.Writer, c.Request)
	})
}
