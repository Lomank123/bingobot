package discord_seeders

import (
	services "bingobot/internal/services"
	"context"
	"crypto/rand"
	"log"
	"math/big"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create Users with Discord Profiles
func GenerateUsersWithProfiles(
	userService *services.UserService,
	profileCollection *mongo.Collection,
	numberOfUsers int,
) []primitive.ObjectID {
	var ids []primitive.ObjectID

	for i := 0; i < numberOfUsers; i++ {
		discordID, err := generateDiscordID()

		if err != nil {
			log.Fatalf("Error generating Discord ID user: %s", err)
		}

		user, err := userService.Create(
			bson.M{
				"discord_id":  discordID,
				"telegram_id": "",
			},
		)

		if err != nil {
			log.Fatalf("Error creating user: %s", err)
		}

		userName, err := generateDiscordUsername()

		if err != nil {
			log.Fatalf("Error generating Discord username user: %s", err)
		}

		// Creating Discord profile
		_, err = profileCollection.InsertOne(
			context.Background(),
			bson.M{
				"user_id":  user.ID,
				"username": userName,
			},
		)

		if err != nil {
			log.Fatalf("Error creating Discord Profile: %s", err)
		}

		ids = append(ids, user.ID)
		log.Printf(
			"User (%s) created with Discord username: %s\n",
			discordID, userName,
		)
	}

	return ids
}

// Generates a random 19-digit string
// Example: 1291447916929744911
func generateDiscordID() (string, error) {
	var id string
	for i := 0; i < 19; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		id += n.String()
	}
	return id, nil
}

// Generates a random 8-character string sequence
// Example: asdfQwRW
func generateDiscordUsername() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result string

	for i := 0; i < 8; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result += string(charset[n.Int64()])
	}

	return result, nil
}
