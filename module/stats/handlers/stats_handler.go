package handlers

import (
	"myapp/common/utils"
	"myapp/module/stats/models"
	"myapp/module/stats/services"

	"github.com/gin-gonic/gin"
)

func GetAllStats(c *gin.Context) {
	page, limit, sortField, sortOrder, _ := utils.ParsePaginationQuery(c)
	stats, total, err := services.GetAllStats(page, limit, sortField, sortOrder)
	if err != nil {
		utils.JSONError(c, err)
		return
	}
	utils.JSONData(c, gin.H{
		"data":        stats,
		"totalCount":  total,
		"totalPages":  utils.TotalPages(total, limit),
		"currentPage": page,
	})
}

func Visit(c *gin.Context) {
	var request models.Stats
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.JSONError(c, err)
		return
	}
	
	if err := services.Visit(request); err != nil {
		utils.JSONError(c, err)
		return
	}

	utils.JSONData(c, gin.H{"message": "Stats updated"})
}

