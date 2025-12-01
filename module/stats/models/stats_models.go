package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Stats struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Ip   string              `bson:"ip" json:"ip"`
	Device string            `bson:"device" json:"device"`
	Area  string             `bson:"area" json:"area"`
	LastTime  string    `bson:"lastTime" json:"lastTime"`
	CountVisit  int     	 `bson:"countVisit" json:"countVisit"`
}

