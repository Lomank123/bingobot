package discord_services

import (
	models "bingobot/internal/models"
	"context"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	Collection *mongo.Collection
}

func (us UserService) GetOrCreateUser(i *discordgo.Interaction) (user *models.User, isCreated bool, error error) {
	discordUser := us.ParseDiscordUser(i)
	user, err := us.FindByDiscordId(discordUser.ID)
	isCreated = false

	if err != nil {
		if err == mongo.ErrNoDocuments {
			user, err = us.CreateUser(discordUser.ID)

			if err != nil {
				return nil, false, err
			}

			isCreated = true
		}
	}

	return user, isCreated, nil
}

func (us UserService) CreateUser(id string) (*models.User, error) {
	newUser := models.User{DiscordID: id}
	_, err := us.Collection.InsertOne(context.Background(), newUser)

	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (us UserService) FindByDiscordId(discordId string) (*models.User, error) {
	result := us.Collection.FindOne(context.Background(), bson.M{"discord_id": discordId})
	var user models.User
	err := result.Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Return discord user instance parsed from interation event
func (UserService) ParseDiscordUser(i *discordgo.Interaction) *discordgo.User {
	if i.Member != nil {
		return i.Member.User
	}

	return i.User
}

func NewUserService(collection *mongo.Collection) *UserService {
	return &UserService{
		Collection: collection,
	}
}
