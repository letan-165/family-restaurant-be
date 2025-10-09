package services

import (
	"errors"
	"myapp/common/errors_code"
	"myapp/common/utils"
	"myapp/config"
	models_item "myapp/module/item/models"
	item_service "myapp/module/item/services"
	"myapp/module/order/models"
	models_order "myapp/module/order/models"
	"myapp/module/order/models/dto"
	user_service "myapp/module/user/services"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// hàm riêng
func buildOrderItems(reqItems []dto.ItemOrderSaveRequest) ([]models_order.ItemOrder, int, error) {
	itemIDs := make([]string, len(reqItems))
	for i, it := range reqItems {
		itemIDs[i] = it.ID
	}

	itemsDB, err := item_service.GetItemsByIDs(itemIDs)
	if err != nil {
		return nil, 0, err
	}

	itemMap := make(map[string]models_item.Item)
	for _, item := range itemsDB {
		itemMap[item.ID.Hex()] = item
	}

	var (
		items      []models_order.ItemOrder
		totalOrder int
	)

	for _, reqItem := range reqItems {
		item, ok := itemMap[reqItem.ID]
		if !ok {
			return nil, 0, errors_code.ITEM_NO_EXISTS
		}

		total := reqItem.Quantity * item.Price
		items = append(items, models_order.ItemOrder{
			Item:     item,
			Quantity: reqItem.Quantity,
			Total:    total,
		})
		totalOrder += total
	}

	return items, totalOrder, nil
}

func GetAllOrder() ([]models_order.Order, error) {
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	cursor, err := config.OrderCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	for cursor.Next(ctx) {
		var order models.Order
		if err := cursor.Decode(&order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func CreateOrder(request dto.OrderSaveRequest) (*primitive.ObjectID, error) {
	var order models_order.Order
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
	items, total, err := buildOrderItems(request.Items)
	if err != nil {
		return nil, err
	}
	order.Items = items
	order.Status = models_order.PENDING
	order.TimeBooking = time.Now()
	order.Total = total

	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	_, err = config.OrderCollection.InsertOne(ctx, order)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, errors_code.ORDER_EXISTS
		}
		return nil, err
	}
	return &order.ID, nil
}

func GetOrderByID(id string) (*models_order.Order, error) {
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors_code.ORDER_NO_EXISTS
	}

	var order models.Order
	err = config.OrderCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors_code.ORDER_NO_EXISTS
		}
		return nil, err
	}

	return &order, nil
}

func UpdateInfoOrder(id string, request dto.OrderSaveRequest) error {
	//check userID (Nếu có)
	if request.Customer.UserID != "" {
		if _, err := user_service.GetUserByID(request.Customer.UserID); err != nil {
			return err
		}
	}

	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	//check order
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors_code.ORDER_NO_EXISTS
	}
	var order models.Order
	err = config.OrderCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors_code.ORDER_NO_EXISTS
		}
		return err
	}

	if order.Status != models.PENDING {
		return errors_code.ORDER_NO_PENDING
	}

	//set order
	items, total, err := buildOrderItems(request.Items)
	if err != nil {
		return err
	}

	res, err := config.OrderCollection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{
			"customer": request.Customer,
			"items":    items,
			"total":    total,
		}},
	)

	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors_code.ORDER_NO_EXISTS
	}

	return nil
}

func UpdatePendingOrder(id string, status string) error {
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	//check order
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors_code.ORDER_NO_EXISTS
	}
	var order models.Order
	err = config.OrderCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors_code.ORDER_NO_EXISTS
		}
		return err
	}

	orderStatus := models.OrderStatus(status)
	if !orderStatus.IsValid() {
		return errors_code.STATUS_ORDER_INVALID
	}

	if order.Status != models.PENDING {
		return errors_code.ORDER_NO_PENDING
	}

	if models.OrderStatus(status) == models.COMPLETED {
		return errors_code.ORDER_NO_CONFIRM
	}

	res, err := config.OrderCollection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{
			"status": status,
		}},
	)

	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors_code.ORDER_NO_EXISTS
	}

	return nil
}

func UpdateConfirmOrder(id string) error {
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	//check order
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors_code.ORDER_NO_EXISTS
	}
	var order models.Order
	err = config.OrderCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors_code.ORDER_NO_EXISTS
		}
		return err
	}
	if order.Status != models.CONFIRMED {
		return errors_code.ORDER_NO_CONFIRM
	}

	res, err := config.OrderCollection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{
			"status":     models.COMPLETED,
			"timeFinish": time.Now(),
		}},
	)

	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors_code.ORDER_NO_EXISTS
	}

	return nil
}
