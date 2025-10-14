package repository

import (
	"myapp/config/db"
	"myapp/module/user/models"
)

type UserRepository struct {
	BaseRepo *db.BaseRepository[models.User]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		BaseRepo: &db.BaseRepository[models.User]{Collection: db.UserCollection},
	}
}
