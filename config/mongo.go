package config

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

	initCollections()
}

func initCollections() {
	ItemCollection = DB.Collection("items")
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := ItemCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}
}
