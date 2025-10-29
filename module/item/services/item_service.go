package services

import (
	"errors"
	"myapp/common/errors_code"
	"myapp/common/utils"
	"myapp/module/item/models"
	"myapp/module/item/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var itemRepo *repository.ItemRepository

func initRepo() {
	if itemRepo == nil {
		itemRepo = repository.NewItemRepository()
	}
}

func GetAllItems(page, limit int, sortField string, sortOrder int) ([]models.Item, int64, error) {
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

	items, err := itemRepo.BaseRepo.FindAll(ctx, map[string]interface{}{}, findOptions)
	if err != nil {
		return nil, 0, err
	}

	total, err := itemRepo.BaseRepo.Count(ctx, map[string]interface{}{})
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func CreateItem(request models.Item) (*primitive.ObjectID, error) {
	initRepo()
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	if !request.Type.IsValid() {
		return nil, errors_code.TYPE_ITEM_INVALID
	}

	request.ID = primitive.NewObjectID()

	res, err := itemRepo.BaseRepo.Insert(ctx, request)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, errors_code.ITEM_EXISTS
		}
		return nil, err
	}

	return res, nil
}

func GetItemByID(id string) (*models.Item, error) {
	initRepo()
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors_code.ITEM_NO_EXISTS
	}

	item, err := itemRepo.BaseRepo.FindByID(ctx, objID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors_code.ITEM_NO_EXISTS
		}
		return nil, err
	}

	return item, nil
}

func UpdateItem(request models.Item) error {
	initRepo()
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	if !request.Type.IsValid() {
		return errors_code.TYPE_ITEM_INVALID
	}

	update := map[string]interface{}{
		"index": request.Index,
		"img":   request.Img,
		"name":  request.Name,
		"type":  request.Type,
		"price": request.Price,
	}

	res, err := itemRepo.BaseRepo.Update(ctx, request.ID, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors_code.ITEM_NO_EXISTS
	}

	return nil
}

func DeleteItem(id string) error {
	initRepo()
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors_code.ITEM_NO_EXISTS
	}

	res, err := itemRepo.BaseRepo.Delete(ctx, objID)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors_code.ITEM_NO_EXISTS
	}

	return nil
}
