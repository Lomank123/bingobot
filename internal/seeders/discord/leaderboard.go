package discord_seeders

import (
	"bingobot/internal/models"
	"context"
	"crypto/rand"
	"log"
	"math/big"
	"time"

	consts "bingobot/internal/consts"
	discord_consts "bingobot/internal/consts/discord"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// For each user insert numberOfRecords records
func GenerateScoresForUsers(userIds []primitive.ObjectID, scoreCollection *mongo.Collection, numberOfRecords int) {
	for _, userId := range userIds {
		for i := 0; i < numberOfRecords; i++ {
			// Generate random score in interval [0, 4]
			score, _ := rand.Int(rand.Reader, big.NewInt(5))
      // TODO: Need to generate different dates for each score record (month, year)
			scoreLog := models.UserScoreRecord{
				UserId:    userId,
				Domain:    consts.DISCORD_DOMAIN,
				Score:     int(score.Int64()) + 1,
				Command:   discord_consts.SEEDER_COMMAND,
				CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
			}
			_, err := scoreCollection.InsertOne(context.Background(), scoreLog)

			if err != nil {
				log.Fatalf("Error creating score record: %s", err)
			}

			log.Printf(
				"Score record created for User (%s) with score: %d\n",
				userId.Hex(),
				scoreLog.Score,
			)
		}
	}
}
