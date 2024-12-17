package main

import (
	config "bingobot/configs/discord"
	"fmt"
	"log"

	consts "bingobot/internal/consts"
	mongodb "bingobot/internal/mongodb"
	seeders "bingobot/internal/seeders/discord"
	services "bingobot/internal/services"

	"github.com/redis/go-redis/v9"
)

func main() {
	log.Println("Start generating users with Discord profiles...")
	config.LoadConfig()

	// TODO: Move initial db/redis/etc setups to separate util which will be used in start scripts
	// DB Setup
	client := mongodb.ConnectToDB(config.Cfg.DBURI)
	database := client.Database(config.Cfg.DBName)

	// TODO: Move to separate method
	userScoreRecordCollection := database.Collection(
		consts.USER_SCORE_RECORD_COLLECTION_NAME,
	)
	scoreService := services.NewScoreService(userScoreRecordCollection)
	scores, err := scoreService.AggregateLeaderboardData()

	if err != nil {
		log.Fatalf("Error aggregating leaderboard data: %s", err)
	}

	// TODO: Remove
	for _, aggregatedScore := range scores {
		fmt.Printf(
			"\nUser (%s) has %d points in %d/%d\n",
			aggregatedScore.UserID,
			aggregatedScore.Score,
			aggregatedScore.Year,
			aggregatedScore.Month,
		)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Cfg.RedisAddress,
		DB:   0,
	})
	leaderboardService := services.NewLeaderboardService(redisClient)
	// TODO: Move recalculation command to util(?) to reuse it in another command
	err = leaderboardService.RecalculateLeaderboard(scores)

	if err != nil {
		log.Fatalf("Error re-calculating leaderboard: %s", err)
	}

	return
	usersCollection := database.Collection(consts.USER_COLLECTION_NAME)
	profileCollection := database.Collection(
		consts.USER_DISCORD_PROFILE_COLLECTION_NAME,
	)
	// userScoreRecordCollection := database.Collection(
	// 	consts.USER_SCORE_RECORD_COLLECTION_NAME,
	// )
	userService := services.NewUserService(usersCollection)

	usersToCreate := 4
	scoreRecordsToCreate := 3

	userIds := seeders.GenerateUsersWithProfiles(
		userService,
		profileCollection,
		usersToCreate,
	)
	seeders.GenerateScoresForUsers(
		userIds,
		userScoreRecordCollection,
		scoreRecordsToCreate,
	)
}
