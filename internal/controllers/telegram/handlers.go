package telegram_handlers

import (
	"fmt"
	"log"

	general_consts "bingobot/internal/consts"
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
		// TODO: Implement command
		response.Text = general_consts.COMMAND_NOT_FOUND_TEXT
	case consts.MY_SCORE_COMMAND:
		// TODO: Implement command
		response.Text = general_consts.COMMAND_NOT_FOUND_TEXT
	case consts.LEADERBOARD_COMMAND:
		// TODO: Implement command
		response.Text = general_consts.COMMAND_NOT_FOUND_TEXT
	default:
		response.Text = general_consts.COMMAND_NOT_FOUND_TEXT
	}

	if _, err := bot.Send(response); err != nil {
		log.Panicf("error occurred during sending response message: %s", err)
	}
}
