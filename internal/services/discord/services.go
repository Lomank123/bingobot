package discord_services

import (
	consts "bingobot/internal/consts"

	"go.mongodb.org/mongo-driver/mongo"

	services "bingobot/internal/services"
)

type DiscordService struct {
	UserService  *UserService
	EchoService  *EchoService
	ScoreService *services.ScoreService
}

func NewDiscordService(database *mongo.Database) *DiscordService {
	usersCollection := database.Collection(consts.USER_COLLECTION_NAME)
	userScoreRecordCollection := database.Collection(
		consts.USER_SCORE_RECORD_COLLECTION_NAME,
	)
	userDiscordProfileCollection := database.Collection(
		consts.USER_DISCORD_PROFILE_COLLECTION_NAME,
	)

	return &DiscordService{
		UserService: NewUserService(
			usersCollection,
			userDiscordProfileCollection,
		),
		EchoService:  NewEchoService(),
		ScoreService: services.NewScoreService(userScoreRecordCollection),
	}
}
