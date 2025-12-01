package services

import (
	"myapp/common/utils"
	"myapp/module/stats/models"
	"myapp/module/stats/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var statsRepo *repository.StatsRepository

func initRepo() {
	if statsRepo == nil {
		statsRepo = repository.NewStatsRepository()
	}
}

func GetAllStats(page, limit int, sortField string, sortOrder int) ([]models.Stats, int64, error) {
	initRepo()
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	findOptions := utils.BuildMongoFindOptions(utils.PaginationOptions{
		Page:      page,
		Limit:     limit,
		SortField: sortField,
		SortOrder: sortOrder,
		Filter:    map[string]interface{}{},
	})

	Stats, err := statsRepo.BaseRepo.FindAll(ctx, map[string]interface{}{}, findOptions)
	if err != nil {
		return nil, 0, err
	}

	total, err := statsRepo.BaseRepo.Count(ctx, map[string]interface{}{})
	if err != nil {
		return nil, 0, err
	}

	return Stats, total, nil
}

func Visit(request models.Stats) (error) {
	initRepo()
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	existingStats, err := statsRepo.FindByIp(ctx, request.Ip)
	if err == nil && existingStats.Ip != "" {
		update := map[string]interface{}{
			"countVisit": existingStats.CountVisit + 1 ,
			"lastTime":   request.LastTime,
		}

		statsRepo.BaseRepo.Update(ctx, existingStats.ID, update)
		return nil
	}

	request.ID = primitive.NewObjectID()
	request.CountVisit = 1
	_, err = statsRepo.BaseRepo.Insert(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

