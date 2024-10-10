package mongo_client

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DBClient *mongo.Client

// Connect to MongoDB and return a Client instance
func ConnectToDB(uri string) error {
	fmt.Println("Connecting to MongoDB...")

	// TODO: Read about contexts
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	DBClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	err = DBClient.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	fmt.Println("Successfully connected to MongoDB!")
	return nil
}
