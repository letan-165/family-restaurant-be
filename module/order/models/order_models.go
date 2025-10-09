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
	UserID   string `bson:"userID" json:"userID"`
	Receiver string `bson:"receiver" json:"receiver"`
	Phone    string `bson:"phone" json:"phone"`
	Address  string `bson:"address" json:"address"`
}

type ItemOrder struct {
	Item     models.Item `bson:"item" json:"item"`
	Quantity int         `bson:"quantity" json:"quantity"`
	Total    int         `bson:"total" json:"total"`
}

type Order struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Customer    CustomerOrder      `bson:"customer" json:"customer"`
	TimeBooking time.Time          `bson:"timeBooking" json:"timeBooking"`
	Status      OrderStatus        `bson:"status" json:"status"`
	TimeFinish  *time.Time         `bson:"timeFinish" json:"timeFinish"`
	Total       int                `bson:"total" json:"total"`
	Items       []ItemOrder        `bson:"items" json:"items"`
}

func (o OrderStatus) IsValid() bool {
	switch o {
	case PENDING, CANCELLED, CONFIRMED, COMPLETED:
		return true
	default:
		return false
	}
}

func IsValid(status string) bool {
	ordeStatus := OrderStatus(status)
	switch ordeStatus {
	case PENDING, CANCELLED, CONFIRMED, COMPLETED:
		return true
	default:
		return false
	}
}
