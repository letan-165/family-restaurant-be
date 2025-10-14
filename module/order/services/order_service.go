package services

import (
	"errors"
	"fmt"
	"myapp/common/errors_code"
	"myapp/common/utils"
	models_item "myapp/module/item/models"
	item_service "myapp/module/item/services"
	models_notification "myapp/module/notification/models"
	notification_service "myapp/module/notification/services"
	models_order "myapp/module/order/models"
	"myapp/module/order/models/dto"
	"myapp/module/order/repository"
	user_service "myapp/module/user/services"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var orderRepo *repository.OrderRepository

func initRepo() {
	if orderRepo == nil {
		orderRepo = repository.NewOrderRepository()
	}
}

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

func GetAllOrders(page, limit int, sortField string, sortOrder int, status string) ([]models_order.Order, int64, error) {
	initRepo()
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}

	findOptions := utils.BuildMongoFindOptions(utils.PaginationOptions{
		Page:      page,
		Limit:     limit,
		SortField: sortField,
		SortOrder: sortOrder,
		Filter:    filter,
	})

	orders, err := orderRepo.BaseRepo.FindAll(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}

	total, err := orderRepo.BaseRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func GetAllOrdersByCustomer(userID string, page, limit int, sortField string, sortOrder int, status string) ([]models_order.Order, int64, error) {
	initRepo()
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	filter := bson.M{"customer.userID": userID}
	if status != "" {
		filter["status"] = status
	}

	findOptions := utils.BuildMongoFindOptions(utils.PaginationOptions{
		Page:      page,
		Limit:     limit,
		SortField: sortField,
		SortOrder: sortOrder,
		Filter:    filter,
	})

	orders, err := orderRepo.BaseRepo.FindAll(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}

	total, err := orderRepo.BaseRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func CreateOrder(request dto.OrderSaveRequest) (*primitive.ObjectID, error) {
	initRepo()
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

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

	res, err := orderRepo.BaseRepo.Insert(ctx, order)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, errors_code.ORDER_EXISTS
		}
		return nil, err
	}

	//Gửi mail (vẫn giữ goroutine như trước)
	go func(o models_order.Order) {
		if err := notification_service.SendMailBooking(o); err != nil {
			fmt.Println("SendMailBooking error:", err)
		}
	}(order)

	//Gửi AlertOrder
	go func(o models_order.Order) {
		alert := models_notification.AlertOrder{
			ID:          primitive.NewObjectID(),
			OrderID:     o.ID.Hex(),
			TimeBooking: o.TimeBooking,
			Message:     fmt.Sprintf("Có đơn hàng mới (Tổng: %d), vui lòng kiểm tra, xác nhận", o.Total),
		}
		models_notification.Notifier.Broadcast(alert)
	}(order)

	return res, nil
}

func GetOrderByID(id string) (*models_order.Order, error) {
	initRepo()
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors_code.ORDER_NO_EXISTS
	}

	order, err := orderRepo.BaseRepo.FindByID(ctx, objID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors_code.ORDER_NO_EXISTS
		}
		return nil, err
	}

	return order, nil
}

func UpdateInfoOrder(id string, request dto.OrderSaveRequest) error {
	initRepo()
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	//check userID (Nếu có)
	if request.Customer.UserID != "" {
		if _, err := user_service.GetUserByID(request.Customer.UserID); err != nil {
			return err
		}
	}

	//check order
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors_code.ORDER_NO_EXISTS
	}

	order, err := orderRepo.BaseRepo.FindByID(ctx, objID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors_code.ORDER_NO_EXISTS
		}
		return err
	}

	if order.Status != models_order.PENDING {
		return errors_code.ORDER_NO_PENDING
	}

	//set order
	items, total, err := buildOrderItems(request.Items)
	if err != nil {
		return err
	}

	res, err := orderRepo.BaseRepo.Update(ctx, objID, bson.M{
		"customer": request.Customer,
		"items":    items,
		"total":    total,
	})
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors_code.ORDER_NO_EXISTS
	}

	return nil
}

func UpdatePendingOrder(id string, status string) error {
	initRepo()
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	//check order
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors_code.ORDER_NO_EXISTS
	}

	order, err := orderRepo.BaseRepo.FindByID(ctx, objID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors_code.ORDER_NO_EXISTS
		}
		return err
	}

	orderStatus := models_order.OrderStatus(status)
	if !orderStatus.IsValid() {
		return errors_code.STATUS_ORDER_INVALID
	}

	if order.Status != models_order.PENDING {
		return errors_code.ORDER_NO_PENDING
	}

	if models_order.OrderStatus(status) == models_order.COMPLETED {
		return errors_code.ORDER_NO_CONFIRM
	}

	res, err := orderRepo.BaseRepo.Update(ctx, objID, bson.M{"status": status})
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors_code.ORDER_NO_EXISTS
	}

	return nil
}

func UpdateConfirmOrder(id string) error {
	initRepo()
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	//check order
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors_code.ORDER_NO_EXISTS
	}

	order, err := orderRepo.BaseRepo.FindByID(ctx, objID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors_code.ORDER_NO_EXISTS
		}
		return err
	}

	if order.Status != models_order.CONFIRMED {
		return errors_code.ORDER_NO_CONFIRM
	}

	res, err := orderRepo.BaseRepo.Update(ctx, objID, bson.M{
		"status":     models_order.COMPLETED,
		"timeFinish": time.Now(),
	})
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors_code.ORDER_NO_EXISTS
	}

	return nil
}
