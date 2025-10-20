package handler

import (
	"myapp/module/notification/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SendRequest struct {
	Target  string `json:"target" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Body    string `json:"body" binding:"required"`
	IsTopic bool   `json:"is_topic"`
}

// Gửi thông báo FCM
func SendNotification(c *gin.Context) {
	var req SendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu gửi không hợp lệ"})
		return
	}

	// Gọi trực tiếp service mà không cần handler struct
	if err := services.SendMessage(req.Target, req.Title, req.Body, req.IsTopic); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Gửi FCM thành công ✅"})
}
