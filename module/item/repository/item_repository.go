package repository

import (
	"context"
	"myapp/config/db"
	models "myapp/module/item/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemRepository struct {
	BaseRepo *db.BaseRepository[models.Item]
}

func NewItemRepository() *ItemRepository {
	return &ItemRepository{
		BaseRepo: &db.BaseRepository[models.Item]{Collection: db.ItemCollection},
	}
}

func (r *ItemRepository) FindItemsByIDs(ctx context.Context, ids []string) ([]models.Item, error) {
	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIDs = append(objectIDs, objID)
	}

	return r.BaseRepo.FindAll(ctx, bson.M{"_id": bson.M{"$in": objectIDs}}, nil)
}
