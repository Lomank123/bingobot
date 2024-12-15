package services

import (
	consts "bingobot/internal/consts"
	discord_consts "bingobot/internal/consts/discord"
	telegram_consts "bingobot/internal/consts/telegram"
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
// ANY action which implies score modification.
// This is the source of truth for the score system.
type ScoreService struct {
	Collection *mongo.Collection
}

// Find score value for a specific command in a given domain
func FindScore(domain string, command string) (int, bool, error) {
	var score int
	var exists bool

	switch domain {
	case consts.DISCORD_DOMAIN:
		score, exists = discord_consts.COMMAND_SCORE_MAPPING[command]
	case consts.TELEGRAM_DOMAIN:
		score, exists = telegram_consts.COMMAND_SCORE_MAPPING[command]
	default:
		return 0, false, fmt.Errorf("domain %s is not valid", domain)
	}

	if !exists {
		return 0, false, fmt.Errorf("command %s is not valid", command)
	}

	return score, exists, nil
}

// Create user score record.
// Return the amount of points added to the user.
func (ss ScoreService) RecordScore(
	user *models.User,
	command string,
	domain string,
) (int, error) {
	score, _, err := FindScore(domain, command)

	if err != nil {
		return 0, err
	}
	if score == 0 {
		return 0, nil
	}

	isLimitExceeded, err := ss.CheckScoreLimitExceeded(user)

	if err != nil {
		return 0, err
	}
	if isLimitExceeded {
		log.Printf("User (%s) has exceeded score limit", user.ID)
		return 0, nil
	}

	scoreLog := models.UserScoreRecord{
		UserId:    user.ID,
		Domain:    domain,
		Score:     score,
		Command:   command,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}
	ss.Collection.InsertOne(context.Background(), scoreLog)

	log.Printf(
		"User (%s) gained %d points from %s for '%s' command",
		user.ID,
		score,
		domain,
		command,
	)

	return score, nil
}

// Check if user has exceeded score limit per amount of time
func (ss ScoreService) CheckScoreLimitExceeded(user *models.User) (bool, error) {
	// Aggregate user score for the last 24 hours
	yesterday := primitive.NewDateTimeFromTime(time.Now().Add(-24 * time.Hour))
	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.D{
				{Key: "user_id", Value: user.ID},
				{Key: "created_at", Value: bson.D{
					{Key: "$gte", Value: yesterday},
				}},
			}},
		},
		{
			{Key: "$group", Value: bson.D{
				// Required by mongo, Value can be nil
				{Key: "_id", Value: "$user_id"},
				{Key: "score", Value: bson.D{
					{Key: "$sum", Value: "$score"}},
				},
			}},
		},
	}
	score, err := ss.aggregateUserScore(pipeline)

	if err != nil {
		return false, err
	}

	return score >= consts.SCORE_LIMIT_PER_DAY, nil
}

// Return aggregated score for a user
func (ss ScoreService) GetUserTotalScore(user *models.User) (int, error) {
	// Find all records by user id and sum all the scores
	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.D{
				{Key: "user_id", Value: user.ID}},
			},
		},
		{
			{Key: "$group", Value: bson.D{
				// Required by mongo, Value can be nil
				{Key: "_id", Value: "$user_id"},
				{Key: "score", Value: bson.D{
					{Key: "$sum", Value: "$score"}},
				},
			}},
		},
	}
	score, err := ss.aggregateUserScore(pipeline)

	if err != nil {
		return 0, err
	}

	return score, nil
}

// TODO: Probably can be refactored
// In pipeline aggregation field must be called 'score'
func (ss ScoreService) aggregateUserScore(pipeline mongo.Pipeline) (int, error) {
	cursor, err := ss.Collection.Aggregate(
		context.Background(),
		pipeline,
	)

	if err != nil {
		return 0, err
	}

	defer cursor.Close(context.Background())

	// Used to unpack mongo result
	var result struct {
		Score int `bson:"score"`
	}

	// Here we seek only first entry
	if cursor.Next(context.Background()) {
		err := cursor.Decode(&result)

		if err != nil {
			return 0, err
		}

		return result.Score, nil
	}

	return 0, nil
}

func NewScoreService(collection *mongo.Collection) *ScoreService {
	return &ScoreService{
		Collection: collection,
	}
}
