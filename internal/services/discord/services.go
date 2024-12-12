package discord_services

import (
	consts "bingobot/internal/consts"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	services "bingobot/internal/services"
)

type DiscordService struct {
	UserService        *services.UserService
	ScoreService       *services.ScoreService
	LeaderboardService *services.LeaderboardService
	ProfileService     *ProfileService
	EchoService        *EchoService
}

func NewDiscordService(
	database *mongo.Database,
	redisClient *redis.Client,
) *DiscordService {
	usersCollection := database.Collection(consts.USER_COLLECTION_NAME)
	userScoreRecordCollection := database.Collection(
		consts.USER_SCORE_RECORD_COLLECTION_NAME,
	)
	userDiscordProfileCollection := database.Collection(
		consts.USER_DISCORD_PROFILE_COLLECTION_NAME,
	)

	return &DiscordService{
		UserService:        services.NewUserService(usersCollection),
		ScoreService:       services.NewScoreService(userScoreRecordCollection),
		LeaderboardService: services.NewLeaderboardService(redisClient),
		ProfileService:     NewProfileService(userDiscordProfileCollection),
		EchoService:        NewEchoService(),
	}
}
