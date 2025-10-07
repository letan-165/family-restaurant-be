package handlers

import (
	"myapp/common/utils"
	"myapp/module/item/models"
	"myapp/module/item/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateItem(c *gin.Context) {
	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		utils.JSONError(c, err)
		return
	}

	id, err := services.CreateItem(item)
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
	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		utils.JSONError(c, err)
		return
	}

	item.ID, _ = primitive.ObjectIDFromHex(c.Param("id"))
	if err := services.UpdateItem(item); err != nil {
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
