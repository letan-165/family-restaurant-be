package routes

import (
	"myapp/module/item/handlers"

	"github.com/gin-gonic/gin"
)

func ItemRoutes(r *gin.Engine) {
	r.POST("/items", handlers.CreateItem)
	r.GET("/items", handlers.GetAllItems)
	r.GET("/items/:id", handlers.GetItemByID)
	r.PUT("/items/:id", handlers.UpdateItem)
	r.DELETE("/items/:id", handlers.DeleteItem)
}
