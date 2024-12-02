package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserScoreProfile struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UserId        primitive.ObjectID `bson:"user_id"`
	DiscordScore  int                `bson:"discord_score"`
	TelegramScore int                `bson:"telegram_score"`
	LastScoreAt   primitive.DateTime `bson:"last_score_at,omitempty"`
}

func (usp *UserScoreProfile) GetTotalScore() int {
	return usp.DiscordScore + usp.TelegramScore
}
