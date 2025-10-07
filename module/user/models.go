package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID      primitive.ObjectID `bson:"_id"`
	Name    string             `bson:"name"`
	Email   string             `bson:"email"`
	Picture string             `bson:"picture"`
	Role    string             `bson:"role"` // ADMIN, CUSTOMER
}
