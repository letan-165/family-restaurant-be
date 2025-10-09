package services

import (
	"errors"
	"myapp/common/errors_code"
	"myapp/common/utils"
	"myapp/config/db"
	"myapp/module/user/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserByID(id string) (*models.User, error) {
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	var user models.User
	err := db.UserCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors_code.USER_NO_EXISTS
		}
		return nil, err
	}

	return &user, nil
}

func CreateUser(request models.User) error {
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	if !request.Role.IsValid() {
		return errors_code.ROLE_USER_INVALID
	}

	_, err := db.UserCollection.InsertOne(ctx, request)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors_code.USER_EXISTS
		}
		return err
	}
	return nil
}
