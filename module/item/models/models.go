package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemType string

const (
	ItemMain  ItemType = "MAIN"
	ItemSide  ItemType = "SIDE"
	ItemDrink ItemType = "DRINK"
)

type Item struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Name  string             `bson:"name" json:"name"`
	Type  ItemType           `bson:"type" json:"type"`
	Price int                `bson:"price" json:"price"`
}

func (t ItemType) IsValid() bool {
	switch t {
	case ItemMain, ItemSide, ItemDrink:
		return true
	default:
		return false
	}
}
