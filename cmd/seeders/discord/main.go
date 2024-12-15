package main

import (
	config "bingobot/configs/discord"
	"log"

	consts "bingobot/internal/consts"
	mongodb "bingobot/internal/mongodb"
	seeders "bingobot/internal/seeders/discord"
	services "bingobot/internal/services"
)

func main() {
	log.Println("Start generating users with Discord profiles...")
	config.LoadConfig()

	// DB Setup
	client := mongodb.ConnectToDB(config.Cfg.DBURI)
	database := client.Database(config.Cfg.DBName)
	usersCollection := database.Collection(consts.USER_COLLECTION_NAME)
	profileCollection := database.Collection(
		consts.USER_DISCORD_PROFILE_COLLECTION_NAME,
	)
	userScoreRecordCollection := database.Collection(
		consts.USER_SCORE_RECORD_COLLECTION_NAME,
	)
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
