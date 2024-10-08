package main

import (
	config "bingobot/configs/telegram"
	"fmt"
	"log"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	config.LoadConfig()
	fmt.Printf("Telegram bot starts... %s", config.Cfg.TelegramBotToken)
	setupBot()
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
