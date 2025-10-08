package services

import (
	"myapp/common/errors_code"
	"myapp/common/utils"
	"myapp/config"
	item_service "myapp/module/item/services"
	"myapp/module/order/models"
	"myapp/module/order/models/dto"
	user_service "myapp/module/user/services"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllOrder() ([]models.Order, error) {
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	cursor, err := config.OrderCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []models.Order
	for cursor.Next(ctx) {
		var item models.Order
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func CreateOrder(request dto.OrderCreateRequest) (*primitive.ObjectID, error) {
	var order models.Order
	order.ID = primitive.NewObjectID()

	//check userID (Nếu có)
	if request.Customer.UserID != "" {
		if _, err := user_service.GetUserByID(request.Customer.UserID); err != nil {
			return nil, err
		}
	}
	//set user
	order.Customer = request.Customer

	//set order
	var totalOrder int
	for _, reqItem := range request.Items {
		item, err := item_service.GetItemByID(reqItem.ID)
		if err != nil {
			return nil, err
		}
		totalItem := reqItem.Quantity * item.Price

		var itemOrder models.ItemOrder
		itemOrder.Item = *item
		itemOrder.Quantity = reqItem.Quantity
		itemOrder.Total = totalItem
		order.Items = append(order.Items, itemOrder)
		totalOrder += totalItem
	}
	order.Status = models.PENDING
	order.TimeBooking = time.Now()
	order.Total = totalOrder

	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	_, err := config.OrderCollection.InsertOne(ctx, order)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, errors_code.ITEM_EXISTS
		}
		return nil, err
	}
	return &order.ID, nil
}
