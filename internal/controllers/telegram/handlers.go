package telegram_handlers

import (
	"fmt"
	"log"

	consts "bingobot/internal/consts/telegram"
	services "bingobot/internal/services/telegram"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Process any single user input from Telegram
func HandleUpdate(
	bot *telegram.BotAPI,
	update *telegram.Update,
	srvs *services.TelegramService,
) {
	// Ignore any non-message updates
	if update.Message == nil {
		return
	}
	// Ignore all non-command messages
	if !update.Message.IsCommand() {
		return
	}

	userId := fmt.Sprintf("%d", update.Message.From.ID)
	user, _, err := srvs.UserService.GetOrCreateUser(userId)

	if err != nil {
		log.Panicf("error occurred during user retrieval/creation: %s", err)
	}

	// 2nd arg is text and by default it's empty
	response := telegram.NewMessage(update.Message.Chat.ID, "")

	// Handle all possible commands
	switch update.Message.Command() {
	case consts.ECHO_COMMAND:
		response.Text = srvs.EchoService.Handle(update, user)
		response.ReplyToMessageID = update.Message.MessageID
	case consts.HELP_COMMAND:
		// TODO: Implement help command
		response.Text = "Help command is not implemented yet."
	default:
		response.Text = "I don't know that command. Try /help to see all available commands."
	}

	if _, err := bot.Send(response); err != nil {
		log.Panicf("error occurred during sending response message: %s", err)
	}
}
