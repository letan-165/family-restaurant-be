package handlers

import (
	"myapp/common/utils"
	"myapp/module/order/models/dto"
	"myapp/module/order/services"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	var request dto.OrderCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.JSONError(c, err)
		return
	}

	id, err := services.CreateOrder(request)
	if err != nil {
		utils.JSONError(c, err)
		return
	}

	utils.JSONData(c, gin.H{"id": id.Hex()})
}
