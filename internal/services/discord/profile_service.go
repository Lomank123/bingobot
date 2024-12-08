package discord_services

import (
	models "bingobot/internal/models"
	"context"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProfileService struct {
	Collection *mongo.Collection
}

func (ps ProfileService) Create(
	user *models.User,
	discordUser *discordgo.User,
) (*models.UserDiscordProfile, error) {
	res, err := ps.Collection.InsertOne(
		context.Background(),
		bson.M{"user_id": user.ID, "username": discordUser.Username},
	)

	if err != nil {
		return nil, err
	}

	profile := models.UserDiscordProfile{
		ID:       res.InsertedID.(primitive.ObjectID),
		UserId:   user.ID,
		Username: discordUser.Username,
	}

	return &profile, nil
}

func NewProfileService(collection *mongo.Collection) *ProfileService {
	return &ProfileService{Collection: collection}
}
