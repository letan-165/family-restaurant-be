package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerOrder struct {
	UserID   string `bson:"userID"`
	Receiver string `bson:"receiver"`
	Phone    string `bson:"phone"`
	Address  string `bson:"address"`
}

type ItemOrder struct {
	Item     primitive.ObjectID `bson:"item"`
	Quantity int                `bson:"quantity"`
	Total    int                `bson:"total"`
}

type Order struct {
	ID            primitive.ObjectID `bson:"_id"`
	Customer      CustomerOrder      `bson:"customer"`
	TimeBooking   time.Time          `bson:"timeBooking"`
	Status        string             `bson:"status"` // PENDING, CANCELLED, CONFIRMED, COMPLETED
	TimeCompleted *time.Time         `bson:"timeCompleted"`
	Total         int                `bson:"total"`
	Items         []ItemOrder        `bson:"items"`
}
