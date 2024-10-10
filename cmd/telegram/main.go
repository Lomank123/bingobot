package main

import (
	config "bingobot/configs/telegram"
	"context"
	"fmt"
	"log"
	"time"

	mongo_client "bingobot/internal/mongodb"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	fmt.Println("Bingo Bot (Telegram) starting...")

	// ENV setup
	config.LoadConfig()

	// DB Setup
	err := mongo_client.ConnectToDB(config.Cfg.DBURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %s", err)
	}

	setupBot()

	// Gracefully disconnect from DB
	// TODO: Test this, if works - add to discord
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer func() {
		if err := mongo_client.DBClient.Disconnect(ctx); err != nil {
			log.Fatalf("Failed to gracefully disconnect from MongoDB: %s", err)
		}
	}()
}

func setupBot() {
	bot, err := telegram.NewBotAPI(config.Cfg.TelegramBotToken)

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account: %s", bot.Self.UserName)

	// Echoes each message as a response to it
	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf(
				"\n\n[%s] %s\n\n",
				update.Message.From.UserName,
				update.Message.Text,
			)

			msg := telegram.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}
