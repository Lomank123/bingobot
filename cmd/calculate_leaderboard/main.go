package main

import (
	config "bingobot/configs/discord"
	consts "bingobot/internal/consts"
	mongodb "bingobot/internal/mongodb"
	services "bingobot/internal/services"
	utils "bingobot/internal/utils/leaderboard"
	"log"

	"github.com/redis/go-redis/v9"
)

// TODO: Get rid of Discord deps
func main() {
	config.LoadConfig()

	// DB Setup
	client := mongodb.ConnectToDB(config.Cfg.DBURI)
	database := client.Database(config.Cfg.DBName)
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Cfg.RedisAddress,
		DB:   0,
	})

	userScoreRecordCollection := database.Collection(
		consts.USER_SCORE_RECORD_COLLECTION_NAME,
	)
	scoreService := services.NewScoreService(userScoreRecordCollection)
	leaderboardService := services.NewLeaderboardService(redisClient)

	err := utils.ResetLeaderboards(scoreService, leaderboardService)

	if err != nil {
		log.Fatalf("could not calculate leaderboards: %s", err)
	}
}
