package repository

import (
	"context"
	"errors"
	"myapp/config/db"
	models "myapp/module/stats/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StatsRepository struct {
	BaseRepo *db.BaseRepository[models.Stats]
}

func NewStatsRepository() *StatsRepository {
	return &StatsRepository{
		BaseRepo: &db.BaseRepository[models.Stats]{Collection: db.StatsCollection},
	}
}

func (r *StatsRepository) FindByIp(ctx context.Context, ip string) (models.Stats, error) {
	var Stats models.Stats

	err := r.BaseRepo.Collection.FindOne(ctx, bson.M{"ip": ip}).Decode(&Stats)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Stats, mongo.ErrNoDocuments
		}
		return Stats, err
	}

	return Stats, nil
}
