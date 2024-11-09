package mongo_client

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect to MongoDB and return a Client instance
func ConnectToDB(uri string) *mongo.Client {
	fmt.Println("Connecting to MongoDB...")

	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(uri),
	)

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %s", err)
	}

	fmt.Println("Successfully connected to MongoDB!")
	return client
}
