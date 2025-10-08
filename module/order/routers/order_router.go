package routers

import (
	"myapp/module/order/handlers"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(r *gin.Engine) {
	r.POST("/orders", handlers.CreateOrder)
}
