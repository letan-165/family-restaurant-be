package routes

import (
	"myapp/module/user/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.GET("/auth/google", handlers.GoogleLogin)
	r.GET("/auth/google/callback", handlers.GoogleCallback)
}
