package main

import (
	config "bingobot/configs/telegram"
	"fmt"
	"log"

	handlers "bingobot/internal/controllers/telegram"
	mongo_client "bingobot/internal/mongodb"
	services "bingobot/internal/services/telegram"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	fmt.Println("Bingo Bot (Telegram) starting...")

	config.LoadConfig()
	client := mongo_client.ConnectToDB(config.Cfg.DBURI)
	database := client.Database(config.Cfg.DBName)
	srvs := services.NewTelegramService(database)
	bot := initBot()

	u := telegram.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	// Iterate eternally over the updates and handle each one
	for update := range updates {
		handlers.HandleUpdate(bot, &update, srvs)
	}
}

func initBot() *telegram.BotAPI {
	bot, err := telegram.NewBotAPI(config.Cfg.TelegramBotToken)

	if err != nil {
		log.Fatalf("Failed to connect to Telegram API: %s", err)
	}

	bot.Debug = config.Cfg.Debug
	log.Printf("Authorized on account: %s", bot.Self.UserName)
	return bot
}
