package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	DiscordID  string             `bson:"discord_id"`
	TelegramID string             `bson:"telegram_id"`
}
