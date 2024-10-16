package mongo_client

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connect to MongoDB and return a Client instance
func ConnectToDB(uri string) (*mongo.Client, error) {
	fmt.Println("Connecting to MongoDB...")

	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(uri),
	)

	if err != nil {
		return nil, err
	}

	// Ping after connection
	err = client.Ping(context.Background(), readpref.Primary())

	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to MongoDB!")
	return client, nil
}
