package repository

import (
	"myapp/config/db"
	models "myapp/module/stats/models"
)

type StatsRepository struct {
	BaseRepo *db.BaseRepository[models.Stats]
}

func NewStatsRepository() *StatsRepository {
	return &StatsRepository{
		BaseRepo: &db.BaseRepository[models.Stats]{Collection: db.StatsCollection},
	}
}