package utils

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PaginationOptions struct {
	Page      int
	Limit     int
	SortField string
	SortOrder int
	Filter    bson.M
}

func BuildMongoFindOptions(p PaginationOptions) *options.FindOptions {
	skip := (p.Page - 1) * p.Limit
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(p.Limit)).
		SetSort(bson.D{{Key: p.SortField, Value: p.SortOrder}})
	return findOptions
}

func TotalPages(total int64, limit int) int {
	return int(math.Ceil(float64(total) / float64(limit)))
}

func ParsePaginationQuery(c *gin.Context) (page int, limit int, sortField string, sortOrder int, status string) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortField = c.DefaultQuery("sortField", "created_at")
	sortOrderStr := c.DefaultQuery("sortOrder", "desc")
	status = c.DefaultQuery("status", "")

	sortOrder = -1
	if sortOrderStr == "asc" {
		sortOrder = 1
	}
	return
}
