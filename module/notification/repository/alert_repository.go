package repository

import (
	"myapp/config/db"
	"myapp/module/notification/models"
)

type AlertRepository struct {
	BaseRepo *db.BaseRepository[models.AlertOrder]
}

func NewAlertRepository() *AlertRepository {
	return &AlertRepository{
		BaseRepo: &db.BaseRepository[models.AlertOrder]{Collection: db.AlertCollection},
	}
}
