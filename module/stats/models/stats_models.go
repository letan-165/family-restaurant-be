package models

type Stats struct {
	ID         string `bson:"_id" json:"id"`
	Ip         string `bson:"ip" json:"ip"`
	Device     string `bson:"device" json:"device"`
	Area       string `bson:"area" json:"area"`
	LastTime   string `bson:"lastTime" json:"lastTime"`
	CountVisit int    `bson:"countVisit" json:"countVisit"`
}
