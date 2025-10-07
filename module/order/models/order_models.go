package models

import (
	"myapp/module/item/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	PENDING   OrderStatus = "PENDING"
	CANCELLED OrderStatus = "CANCELLED"
	CONFIRMED OrderStatus = "CONFIRMED"
	COMPLETED OrderStatus = "COMPLETED"
)

type CustomerOrder struct {
	UserID   string `bson:"userID"`
	Receiver string `bson:"receiver"`
	Phone    string `bson:"phone"`
	Address  string `bson:"address"`
}

type ItemOrder struct {
	Item     models.Item `bson:"item"`
	Quantity int         `bson:"quantity"`
	Total    int         `bson:"total"`
}

type Order struct {
	ID            primitive.ObjectID `bson:"_id"`
	Customer      CustomerOrder      `bson:"customer"`
	TimeBooking   time.Time          `bson:"timeBooking"`
	Status        OrderStatus        `bson:"status"`
	TimeCompleted *time.Time         `bson:"timeCompleted"`
	Total         int                `bson:"total"`
	Items         []ItemOrder        `bson:"items"`
}

func (o OrderStatus) IsValid() bool {
	switch o {
	case PENDING, CANCELLED, CONFIRMED, COMPLETED:
		return true
	default:
		return false
	}
}
