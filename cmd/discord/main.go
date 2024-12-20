package main

import (
	config "bingobot/configs/discord"
	"fmt"
	"log"
	"os"
	"os/signal"

	handlers "bingobot/internal/controllers/discord"
	mongodb "bingobot/internal/mongodb"
	services "bingobot/internal/services/discord"

	"github.com/bwmarrin/discordgo"
	"github.com/redis/go-redis/v9"
)

func main() {
	fmt.Println("Bingo Bot (Discord) starting...")
	config.LoadConfig()

	// DB Setup
	client := mongodb.ConnectToDB(config.Cfg.DBURI)
	database := client.Database(config.Cfg.DBName)

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Cfg.RedisAddress,
		DB:   0, // Use default DB
	})
	services := services.NewDiscordService(database, redisClient)
	session := initBot(services)
	startServing(session)
}

func initBot(services *services.DiscordService) *discordgo.Session {
	session, _ := discordgo.New("Bot " + config.Cfg.DiscordBotToken)

	// Handlers
	handlers.SetupHandlers(session, services)

	// Commands
	_, err := session.ApplicationCommandBulkOverwrite(
		config.Cfg.DiscordAppId,
		"",
		config.Commands,
	)

	if err != nil {
		log.Fatalf("couldn't register commands: %s", err)
	}

	return session
}

func startServing(s *discordgo.Session) {
	// Lifecycle?
	// TODO: Check what this does exactly
	err := s.Open()

	if err != nil {
		log.Fatalf("could not open session: %s", err)
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	err = s.Close()

	if err != nil {
		log.Printf("couldn't close session gracefully: %s", err)
	}
}
