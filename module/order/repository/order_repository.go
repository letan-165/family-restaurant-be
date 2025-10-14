package repository

import (
	"myapp/config/db"
	"myapp/module/order/models"
)

type OrderRepository struct {
	BaseRepo *db.BaseRepository[models.Order]
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		BaseRepo: &db.BaseRepository[models.Order]{Collection: db.OrderCollection},
	}
}
