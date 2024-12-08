package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Represents result of a single action made by user.
// Action can be anything that implies score modification.
type UserScoreRecord struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserId    primitive.ObjectID `bson:"user_id"`
	Domain    string             `bson:"domain"`     // To distinguish scores from different domains
	Score     int                `bson:"score"`      // To count total score for user
	Command   string             `bson:"command"`    // To distinguish commands for user stats
	CreatedAt primitive.DateTime `bson:"created_at"` // To check score limit per amount of time
}
