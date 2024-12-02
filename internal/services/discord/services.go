package discord_services

import (
	consts "bingobot/internal/consts"

	"go.mongodb.org/mongo-driver/mongo"
)

type DiscordService struct {
	UserService  *UserService
	EchoService  *EchoService
	ScoreService *ScoreService
}

func NewDiscordService(database *mongo.Database) *DiscordService {
	usersCollection := database.Collection(consts.USER_COLLECTION_NAME)
	userScoreProfileCollection := database.Collection(
		consts.USER_SCORE_PROFILE_COLLECTION_NAME,
	)
	userDiscordProfileCollection := database.Collection(
		consts.USER_DISCORD_PROFILE_COLLECTION_NAME,
	)

	return &DiscordService{
		UserService: NewUserService(
			usersCollection,
			userDiscordProfileCollection,
			userScoreProfileCollection,
		),
		EchoService:  NewEchoService(),
		ScoreService: NewScoreService(userScoreProfileCollection),
	}
}
