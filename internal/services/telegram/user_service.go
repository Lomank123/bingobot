package telegram_services

import (
	models "bingobot/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	Collection *mongo.Collection
}

func (us UserService) GetOrCreateUser(id string) (user *models.User, isCreated bool, error error) {
	user, err := us.FindByTelegramId(id)
	isCreated = false

	if err != nil {
		if err == mongo.ErrNoDocuments {
			user, err = us.CreateUser(id)

			if err != nil {
				return nil, false, err
			}

			isCreated = true
		}
	}

	return user, isCreated, nil
}

func (us UserService) CreateUser(telegramId string) (*models.User, error) {
	newUser := models.User{TelegramID: telegramId}
	_, err := us.Collection.InsertOne(context.Background(), newUser)

	if err != nil {
		return nil, err
	}

	// TODO: Create telegram profile and score profile here

	return &newUser, nil
}

func (us UserService) FindByTelegramId(telegramId string) (*models.User, error) {
	result := us.Collection.FindOne(context.Background(), bson.M{"telegram_id": telegramId})
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
