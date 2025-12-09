package database

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	clientOnce sync.Once
	dbName     string
)

// init loads environment vars once at package initialization

func init() {
	_ = godotenv.Load(".env")
	// if err != nil {
	// 	panic(err)
	// }
	dbName = os.Getenv("DATABASE_NAME")
	if dbName == "" {
		log.Fatal("DATABASE_NAME not set in environment")
	}
	log.Println("init of database file loaded")
}

func GetClient() *mongo.Client {
	clientOnce.Do(func() {
		uri := os.Getenv("MONGODB_URI")
		if uri == "" {
			log.Fatal("MONGODB_URI not set in environment")
		}

		opts := options.Client().ApplyURI(uri).
			SetConnectTimeout(time.Second * 5).         // Prevents app hanging if Mongo is unreachable
			SetServerSelectionTimeout(5 * time.Second). // Determines how long the driver waits to find a healthy node
			SetRetryWrites(true).                       // MongoDB standard — automatic retry of safe operations
			SetMaxPoolSize(20).                         // Worker pool size — controls concurrency
			SetMinPoolSize(5).                          // Keeps a warm pool of connections
			SetMaxConnIdleTime(30 * time.Second)        // Ensures stale connections are cleaned up

		c, err := mongo.Connect(context.Background(), opts)

		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}

		if err := c.Ping(context.Background(), nil); err != nil {
			log.Fatalf("MongoDB ping failed: %v", err)
		}

		client = c
		log.Println("MongoDB connected successfully")
	})

	return client
}

func OpenCollection(name string) *mongo.Collection {
	return GetClient().Database(dbName).Collection(name)
}

func Disconnect() {
	if client != nil {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		} else {
			log.Println("MongoDB disconnected successfully")
		}
	}
}
