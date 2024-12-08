package services

import (
	models "bingobot/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Operates high-level user-related actions
type UserService struct {
	Collection *mongo.Collection
}

func (us UserService) GetOrCreate(discordId string, telegramId string) (
	user *models.User,
	isCreated bool,
	error error,
) {
	filters := bson.M{}

	if discordId != "" {
		filters["discord_id"] = discordId
	}
	if telegramId != "" {
		filters["telegram_id"] = telegramId
	}

	user, err := us.FindBy(filters)
	isCreated = false

	if err != nil {
		if err == mongo.ErrNoDocuments {
			userData := bson.M{
				"discord_id":  discordId,
				"telegram_id": telegramId,
			}
			user, err = us.Create(userData)

			if err != nil {
				return nil, false, err
			}

			isCreated = true
		}
	}

	return user, isCreated, nil
}

func (us UserService) Create(userData bson.M) (*models.User, error) {
	res, err := us.Collection.InsertOne(context.Background(), userData)

	if err != nil {
		return nil, err
	}

	user := models.User{
		ID:         res.InsertedID.(primitive.ObjectID),
		DiscordID:  userData["discord_id"].(string),
		TelegramID: userData["telegram_id"].(string),
	}

	return &user, nil
}

func (us UserService) FindBy(filters bson.M) (*models.User, error) {
	result := us.Collection.FindOne(context.Background(), filters)
	var user models.User
	err := result.Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func NewUserService(collection *mongo.Collection) *UserService {
	return &UserService{
		Collection: collection,
	}
}
