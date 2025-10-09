package handlers

import (
	"myapp/common/utils"
	"myapp/module/order/models/dto"
	"myapp/module/order/services"

	"github.com/gin-gonic/gin"
)

func GetAllOrders(c *gin.Context) {
	orders, err := services.GetAllOrder()
	if err != nil {
		utils.JSONError(c, err)
		return
	}
	utils.JSONData(c, orders)
}

func CreateOrder(c *gin.Context) {
	var request dto.OrderSaveRequest
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

func GetOrderByID(c *gin.Context) {
	item, err := services.GetOrderByID(c.Param("id"))
	if err != nil {
		utils.JSONError(c, err)
		return
	}
	utils.JSONData(c, item)
}

func UpdateInfoOrder(c *gin.Context) {
	var request dto.OrderSaveRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.JSONError(c, err)
		return
	}

	err := services.UpdateInfoOrder(c.Param("id"), request)
	if err != nil {
		utils.JSONError(c, err)
		return
	}

	utils.JSONData(c, gin.H{"message": "Update success"})
}

func UpdatePendingOrder(c *gin.Context) {
	err := services.UpdatePendingOrder(c.Param("id"), c.Param("status"))
	if err != nil {
		utils.JSONError(c, err)
		return
	}

	utils.JSONData(c, gin.H{"message": "Update success"})
}

func UpdateConfirmOrder(c *gin.Context) {
	err := services.UpdateConfirmOrder(c.Param("id"))
	if err != nil {
		utils.JSONError(c, err)
		return
	}

	utils.JSONData(c, gin.H{"message": "Update success"})
}
