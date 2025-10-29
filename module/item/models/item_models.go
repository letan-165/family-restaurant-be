package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemType string

const (
	MAIN  ItemType = "MAIN"
	SIDE  ItemType = "SIDE"
	DRINK ItemType = "DRINK"
)

type Item struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Img   string             `bson:"img" json:"img"`
	Index int                `bson:"index" json:"index"`
	Name  string             `bson:"name" json:"name"`
	Type  ItemType           `bson:"type" json:"type"`
	Price int                `bson:"price" json:"price"`
}

func (t ItemType) IsValid() bool {
	switch t {
	case MAIN, SIDE, DRINK:
		return true
	default:
		return false
	}
}
