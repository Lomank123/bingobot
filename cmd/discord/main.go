package main

import (
	config "bingobot/configs/discord"
	"fmt"
	"log"
	"os"
	"os/signal"

	handlers "bingobot/internal/controllers/discord"
	mongo_client "bingobot/internal/mongodb"

	"github.com/bwmarrin/discordgo"
)

func main() {
	fmt.Println("Bingo Bot (Discord) starting...")
	config.LoadConfig()

	// DB Setup
	err := mongo_client.ConnectToDB(config.Cfg.DBURI, config.Cfg.DBName)

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %s", err)
	}

	session := initBot()
	startServing(session)
}

func initBot() *discordgo.Session {
	session, _ := discordgo.New("Bot " + config.Cfg.DiscordBotToken)

	// TODO: Pass initialized services to handlers
	// Handlers
	handlers.SetupHandlers(session)

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
