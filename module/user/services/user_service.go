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

func CreateOrGetUser(request models.User) (models.User, error) {
	initRepo()
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	existingUser, err := userRepo.FindByEmail(ctx, request.Email)
	if err == nil && existingUser.Email != "" {
		return existingUser, nil
	}

	request.Role = models.CUSTOMER

	_, err = userRepo.BaseRepo.Insert(ctx, request)
	if err != nil {
		return request, err
	}

	return request, nil
}
