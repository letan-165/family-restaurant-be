package db

import (
	"context"
	"log"
	"myapp/common/utils"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var ItemCollection *mongo.Collection
var OrderCollection *mongo.Collection
var UserCollection *mongo.Collection
var StatsCollection *mongo.Collection

func ConnectMongo() {
	godotenv.Load()
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")

	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := utils.DefaultCtx()
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	DB = client.Database(dbName)
	log.Println("Connected to MongoDB Atlas")

	InitCollections()
}
func InitCollections() {
	ItemCollection = DB.Collection("items")
	OrderCollection = DB.Collection("orders")
	UserCollection = DB.Collection("users")
	StatsCollection = DB.Collection("stats")


	indexes := []struct {
		coll   *mongo.Collection
		field  string
		unique bool
	}{
		{ItemCollection, "name", true},
		{ItemCollection, "index", true},
		{UserCollection, "email", true},
	}

	for _, idx := range indexes {
		model := mongo.IndexModel{
			Keys:    bson.M{idx.field: 1},
			Options: options.Index().SetUnique(idx.unique),
		}
		if _, err := idx.coll.Indexes().CreateOne(context.TODO(), model); err != nil {
			log.Fatalf("Không thể tạo index cho %s.%s: %v", idx.coll.Name(), idx.field, err)
		}
	}

	log.Println("Đã khởi tạo collections và indexes thành công!")
}
