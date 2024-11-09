package telegram_services

import (
	consts "bingobot/internal/consts"

	"go.mongodb.org/mongo-driver/mongo"
)

type TelegramService struct {
	UserService *UserService
	EchoService *EchoService
}

func NewTelegramService(database *mongo.Database) *TelegramService {
	usersCollection := database.Collection(consts.USER_COLLECTION_NAME)

	return &TelegramService{
		UserService: NewUserService(usersCollection),
		EchoService: NewEchoService(),
	}
}
