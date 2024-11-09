package telegram_services

import (
	"bingobot/internal/models"
	"fmt"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type EchoService struct{}

func (EchoService) Handle(
	update *telegram.Update,
	user *models.User,
) string {
	// TODO: Remove command part of the message (e.g. `/echo `).
	//  Also replace user id with username
	return fmt.Sprintf("**%s says:** %s", user.TelegramID, update.Message.Text)
}

func NewEchoService() *EchoService {
	return &EchoService{}
}
