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

func (us UserService) GetOrCreateUser(i *discordgo.Interaction) *models.User {
	discordUser := us.ParseDiscordUser(i)
	user, err := us.FindByDiscordId(discordUser.ID)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			user = us.CreateUser(discordUser.ID)
		}
	}

	return user
}

func (us UserService) CreateUser(id string) *models.User {
	newUser := models.User{DiscordID: id}
	_, err := us.Collection.InsertOne(context.Background(), newUser)

	if err != nil {
		panic(err)
	}

	return &newUser
}

// Find user by discord user id
func (us UserService) FindByDiscordId(discordId string) (*models.User, error) {
	result := us.Collection.FindOne(context.Background(), bson.M{"discord_id": discordId})

	var user models.User
	err := result.Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

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
