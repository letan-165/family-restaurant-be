package handler

import (
	"myapp/common/utils"
	"myapp/module/notification/services"
	"myapp/module/order/models"

	"github.com/gin-gonic/gin"
)

type SendRequest struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}

func SendNotification(c *gin.Context) {
	var req models.Order
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, err)
		return
	}

	if err := services.SendFCMBooking(req); err != nil {
		utils.JSONError(c, err)
		return
	}
	utils.JSONData(c, gin.H{"message": "Gửi FCM thành công"})
}
