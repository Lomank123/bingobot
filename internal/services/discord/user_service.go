package discord_services

import (
	models "bingobot/internal/models"
	"context"
	"log"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	Collection                   *mongo.Collection
	UserDiscordProfileCollection *mongo.Collection
	UserScoreCollection          *mongo.Collection
}

func (us UserService) GetOrCreateUser(i *discordgo.Interaction) (user *models.User, isCreated bool, error error) {
	discordUser := us.ParseDiscordUser(i)
	user, err := us.FindByDiscordId(discordUser.ID)
	isCreated = false

	if err != nil {
		if err == mongo.ErrNoDocuments {
			user, err = us.CreateUser(discordUser)

			if err != nil {
				return nil, false, err
			}

			isCreated = true
		}
	}

	return user, isCreated, nil
}

func (us UserService) CreateUser(discordUser *discordgo.User) (*models.User, error) {
	newUser := models.User{DiscordID: discordUser.ID}
	res, err := us.Collection.InsertOne(context.Background(), newUser)

	if err != nil {
		return nil, err
	}

	// Otherwise it will be set to ObjectId('000...00')
	newUser.ID = res.InsertedID.(primitive.ObjectID)

	_, err = us.UserDiscordProfileCollection.InsertOne(
		context.Background(),
		bson.M{"user_id": newUser.ID, "username": discordUser.Username},
	)

	// TODO: Add error handling
	// TODO: Make all db insertions in a transaction
	if err != nil {
		log.Panicf("could not insert discord profile for user %s: %s", newUser.ID, err)
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

func NewUserService(
	collection *mongo.Collection,
	userDiscordProfileCollection *mongo.Collection,
) *UserService {
	return &UserService{
		Collection:                   collection,
		UserDiscordProfileCollection: userDiscordProfileCollection,
	}
}
