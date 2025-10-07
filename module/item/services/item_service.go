package services

import (
	"myapp/common/utils"
	"myapp/config"
	models "myapp/module/item/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var itemCollection *mongo.Collection

func InitCollections() {
	itemCollection = config.DB.Collection("items")
}

func CreateItem(item models.Item) (*primitive.ObjectID, error) {
	item.ID = primitive.NewObjectID()
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	_, err := itemCollection.InsertOne(ctx, item)
	if err != nil {
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

	objID, _ := primitive.ObjectIDFromHex(id)
	var item models.Item
	err := itemCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func UpdateItem(item models.Item) error {
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	_, err := itemCollection.UpdateOne(
		ctx,
		bson.M{"_id": item.ID},
		bson.M{"$set": bson.M{
			"name":  item.Name,
			"type":  item.Type,
			"price": item.Price,
		}},
	)
	return err
}

func DeleteItem(id string) error {
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(id)
	_, err := itemCollection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
