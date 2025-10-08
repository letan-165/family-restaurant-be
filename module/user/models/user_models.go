package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole string

const (
	ADMIN    UserRole = "ADMIN"
	CUSTOMER UserRole = "CUSTOMER"
)

type User struct {
	ID      primitive.ObjectID `bson:"_id"`
	Name    string             `bson:"name"`
	Email   string             `bson:"email"`
	Picture string             `bson:"picture"`
	Role    UserRole           `bson:"role"`
}
