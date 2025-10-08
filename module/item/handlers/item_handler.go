package handlers

import (
	"myapp/common/utils"
	"myapp/module/item/models"
	"myapp/module/item/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateItem(c *gin.Context) {
	var request models.Item
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.JSONError(c, err)
		return
	}

	id, err := services.CreateItem(request)
	if err != nil {
		utils.JSONError(c, err)
		return
	}

	utils.JSONData(c, gin.H{"id": id.Hex()})
}

func GetAllItems(c *gin.Context) {
	items, err := services.GetAllItems()
	if err != nil {
		utils.JSONError(c, err)
		return
	}
	utils.JSONData(c, items)
}

func GetItemByID(c *gin.Context) {
	item, err := services.GetItemByID(c.Param("id"))
	if err != nil {
		utils.JSONError(c, err)
		return
	}
	utils.JSONData(c, item)
}

func UpdateItem(c *gin.Context) {
	var request models.Item
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.JSONError(c, err)
		return
	}

	request.ID, _ = primitive.ObjectIDFromHex(c.Param("id"))
	if err := services.UpdateItem(request); err != nil {
		utils.JSONError(c, err)
		return
	}

	utils.JSONData(c, gin.H{"message": "Item updated"})
}

func DeleteItem(c *gin.Context) {
	if err := services.DeleteItem(c.Param("id")); err != nil {
		utils.JSONError(c, err)
		return
	}
	utils.JSONData(c, gin.H{"message": "Item deleted"})
}
