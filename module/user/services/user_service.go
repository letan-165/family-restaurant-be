package services

import (
	"errors"
	"myapp/common/errors_code"
	"myapp/common/utils"
	"myapp/module/user/models"
	"myapp/module/user/repository"

	"go.mongodb.org/mongo-driver/mongo"
)

var userRepo *repository.UserRepository

func initRepo() {
	if userRepo == nil {
		userRepo = repository.NewUserRepository()
	}
}

func GetUserByID(id string) (*models.User, error) {
	initRepo()
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	user, err := userRepo.BaseRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors_code.USER_NO_EXISTS
		}
		return nil, err
	}

	return user, nil
}

func CreateUser(request models.User) error {
	initRepo()
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	if !request.Role.IsValid() {
		return errors_code.ROLE_USER_INVALID
	}

	_, err := userRepo.BaseRepo.Insert(ctx, request)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors_code.USER_EXISTS
		}
		return err
	}
	return nil
}
