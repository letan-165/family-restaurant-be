package routers

import (
	"myapp/module/notification/handler"

	"github.com/gin-gonic/gin"
)

func FCMRouter(r *gin.Engine) {
	r.POST("/send", handler.SendNotification)
}
