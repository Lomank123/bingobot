package discord_services

import (
	models "bingobot/internal/models"
	"context"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DiscordUserService struct {
	Collection *mongo.Collection
}

func (dus DiscordUserService) GetOrCreateUser(i *discordgo.Interaction) *models.User {
	discordUser := dus.ParseDiscordUser(i)
	user, err := dus.FindByDiscordId(discordUser.ID)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			user = dus.CreateUser(discordUser.ID)
		}
	}

	return user
}

func (dus DiscordUserService) CreateUser(id string) *models.User {
	newUser := models.User{DiscordID: id}
	_, err := dus.Collection.InsertOne(context.Background(), newUser)

	if err != nil {
		panic(err)
	}

	return &newUser
}

// Find user by discord user id
func (dus DiscordUserService) FindByDiscordId(discordId string) (*models.User, error) {
	result := dus.Collection.FindOne(context.Background(), bson.M{"discord_id": discordId})

	var user models.User
	err := result.Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (DiscordUserService) ParseDiscordUser(i *discordgo.Interaction) *discordgo.User {
	if i.Member != nil {
		return i.Member.User
	}

	return i.User
}
