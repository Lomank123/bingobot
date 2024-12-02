package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserTelegramProfile struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	UserId   primitive.ObjectID `bson:"user_id"`
	Username string             `bson:"username"`
}
