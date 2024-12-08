package discord_services

import (
	consts "bingobot/internal/consts"

	"go.mongodb.org/mongo-driver/mongo"

	services "bingobot/internal/services"
)

type DiscordService struct {
	UserService    *services.UserService
	ScoreService   *services.ScoreService
	ProfileService *ProfileService
	EchoService    *EchoService
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
		UserService:    services.NewUserService(usersCollection),
		ScoreService:   services.NewScoreService(userScoreRecordCollection),
		ProfileService: NewProfileService(userDiscordProfileCollection),
		EchoService:    NewEchoService(),
	}
}
