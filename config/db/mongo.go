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

func ConnectMongo() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

	itemIndex := mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	}
	if _, err := ItemCollection.Indexes().CreateOne(context.TODO(), itemIndex); err != nil {
		log.Fatalf("Không thể tạo index cho items.name: %v", err)
	}

	userIndex := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}
	if _, err := UserCollection.Indexes().CreateOne(context.TODO(), userIndex); err != nil {
		log.Fatalf("Không thể tạo index cho users.email: %v", err)
	}

	log.Println("Đã khởi tạo collections và indexes thành công!")
}
