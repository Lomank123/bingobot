package discord_services

import (
	"fmt"

	"bingobot/internal/models"
	utils "bingobot/internal/utils/discord"
)

type EchoService struct{}

func (EchoService) Handle(
	opts utils.OptionMap,
	user *models.User,
) string {
	// Pass user data to service function
	message := opts["message"].StringValue()
	return fmt.Sprintf("**%s says:** %s", user.DiscordID, message)
}

func NewEchoService() *EchoService {
	return &EchoService{}
}
