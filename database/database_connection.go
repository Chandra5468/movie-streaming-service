package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBInstance() *mongo.Client {
	err := godotenv.Load(".env")

	if err != nil {
		log.Printf("warning: unabel to find .env file %v\n", err)
	}

	MongoDB := os.Getenv("MONGODB_URI")

	if MongoDB == "" {
		log.Fatal("mongo uri not configured")
	}

	clientOptions := options.Client().ApplyURI(MongoDB)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil
	}

	return client
}

var Client *mongo.Client = DBInstance()

func OpenCollection(collectionName string) *mongo.Collection {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("unabel to find .env file %v \n", err)
	}

	databaseName := os.Getenv("DATABASE_NAME")

	collection := Client.Database(databaseName).Collection(collectionName)

	return collection
}
