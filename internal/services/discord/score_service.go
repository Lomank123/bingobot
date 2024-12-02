package discord_services

import (
	consts "bingobot/internal/consts/discord"
	"bingobot/internal/models"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Score service is responsible for
// ANY action which implies score modification
type ScoreService struct {
	Collection *mongo.Collection
}

func (ss ScoreService) IncrementScore(user *models.User, command string) error {
	score, exists := consts.COMMAND_SCORE_MAPPING[command]

	if !exists {
		return fmt.Errorf("command %s is not valid", command)
	}

	// TODO: Add barrier for score incrementation
	// TODO: Change message. Add command name and user data other than ID
	log.Printf("Incrementing score for user %s by %d", user.ID, score)

	// Adding score and updating last_score_at field
	ss.Collection.FindOneAndUpdate(
		context.Background(),
		bson.M{"user_id": user.ID},
		bson.M{
			"$inc": bson.M{"discord_score": score},
			"$set": bson.M{"last_score_at": primitive.NewDateTimeFromTime(time.Now())},
		},
	)

	return nil
}

func (ss ScoreService) GetTotalUserScore(user *models.User) (int, error) {
	result := ss.Collection.FindOne(
		context.Background(),
		bson.M{"user_id": user.ID},
	)
	var userScoreProfile models.UserScoreProfile
	err := result.Decode(&userScoreProfile)

	if err != nil {
		return 0, err
	}

	return userScoreProfile.GetTotalScore(), nil
}

func NewScoreService(collection *mongo.Collection) *ScoreService {
	return &ScoreService{
		Collection: collection,
	}
}
