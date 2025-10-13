package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AlertOrder struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	OrderID     string             `bson:"orderID" json:"orderID"`
	TimeBooking time.Time          `bson:"timeBooking" json:"timeBooking"`
	Message     string             `bson:"message" json:"message"`
}
