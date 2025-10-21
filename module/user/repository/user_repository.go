package repository

import (
	"context"
	"errors"
	"myapp/config/db"
	"myapp/module/user/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	BaseRepo *db.BaseRepository[models.User]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		BaseRepo: &db.BaseRepository[models.User]{Collection: db.UserCollection},
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	err := r.BaseRepo.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, mongo.ErrNoDocuments
		}
		return user, err
	}

	return user, nil
}
