package discord_services

import (
	consts "bingobot/internal/consts/discord"

	"go.mongodb.org/mongo-driver/mongo"
)

type DiscordService struct {
	UserService *UserService
	EchoService *EchoService
}

func NewDiscordService(database *mongo.Database) *DiscordService {
	usersCollection := database.Collection(consts.USER_COLLECTION_NAME)

	return &DiscordService{
		UserService: NewUserService(usersCollection),
		EchoService: NewEchoService(),
	}
}
