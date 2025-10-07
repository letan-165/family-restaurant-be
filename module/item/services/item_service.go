package services

import (
	"context"
	"errors"
	"log"
	"myapp/common/errors_code"
	"myapp/common/utils"
	"myapp/config"
	models "myapp/module/item/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var itemCollection = config.ItemCollection

func InitCollections() {
	itemCollection = config.DB.Collection("items")
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := itemCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}
}

func CreateItem(item models.Item) (*primitive.ObjectID, error) {
	item.ID = primitive.NewObjectID()
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	if !item.Type.IsValid() {
		return nil, errors_code.TYPE_ITEM_INVALID
	}

	_, err := itemCollection.InsertOne(ctx, item)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, errors_code.ITEM_EXISTS
		}
		return nil, err
	}
	return &item.ID, nil
}

func GetAllItems() ([]models.Item, error) {
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	cursor, err := itemCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []models.Item
	for cursor.Next(ctx) {
		var item models.Item
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func GetItemByID(id string) (*models.Item, error) {
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var item models.Item
	err = itemCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors_code.ITEM_NO_EXISTS
		}
		return nil, err
	}

	return &item, nil
}

func UpdateItem(item models.Item) error {
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	if !item.Type.IsValid() {
		return errors_code.TYPE_ITEM_INVALID
	}

	res, err := itemCollection.UpdateOne(
		ctx,
		bson.M{"_id": item.ID},
		bson.M{"$set": bson.M{
			"name":  item.Name,
			"type":  item.Type,
			"price": item.Price,
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
	res, err := itemCollection.DeleteOne(ctx, bson.M{"_id": objID})

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors_code.ITEM_NO_EXISTS
	}

	return nil
}
