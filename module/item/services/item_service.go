package services

import (
	"errors"
	"myapp/common/errors_code"
	"myapp/common/utils"
	"myapp/config/db"
	models "myapp/module/item/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllItems(page, limit int, sortField string, sortOrder int) ([]models.Item, int64, error) {
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	findOptions := utils.BuildMongoFindOptions(utils.PaginationOptions{
		Page:      page,
		Limit:     limit,
		SortField: sortField,
		SortOrder: sortOrder,
		Filter:    bson.M{},
	})

	cursor, err := db.ItemCollection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var items []models.Item
	for cursor.Next(ctx) {
		var item models.Item
		if err := cursor.Decode(&item); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}

	total, err := db.OrderCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func CreateItem(request models.Item) (*primitive.ObjectID, error) {
	request.ID = primitive.NewObjectID()
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	if !request.Type.IsValid() {
		return nil, errors_code.TYPE_ITEM_INVALID
	}

	_, err := db.ItemCollection.InsertOne(ctx, request)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, errors_code.ITEM_EXISTS
		}
		return nil, err
	}
	return &request.ID, nil
}

func GetItemByID(id string) (*models.Item, error) {
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors_code.ITEM_NO_EXISTS
	}

	var item models.Item
	err = db.ItemCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors_code.ITEM_NO_EXISTS
		}
		return nil, err
	}

	return &item, nil
}

func UpdateItem(request models.Item) error {
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	if !request.Type.IsValid() {
		return errors_code.TYPE_ITEM_INVALID
	}

	res, err := db.ItemCollection.UpdateOne(
		ctx,
		bson.M{"_id": request.ID},
		bson.M{"$set": bson.M{
			"name":  request.Name,
			"type":  request.Type,
			"price": request.Price,
		}},
	)

	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors_code.ITEM_NO_EXISTS
	}

	return nil
}

func DeleteItem(id string) error {
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(id)
	res, err := db.ItemCollection.DeleteOne(ctx, bson.M{"_id": objID})

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors_code.ITEM_NO_EXISTS
	}

	return nil
}

func GetItemsByIDs(ids []string) ([]models.Item, error) {
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors_code.ITEM_NO_EXISTS
		}
		objectIDs = append(objectIDs, objID)
	}

	filter := bson.M{"_id": bson.M{"$in": objectIDs}}

	cursor, err := db.ItemCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []models.Item
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, errors_code.ITEM_NO_EXISTS
	}

	return items, nil
}
