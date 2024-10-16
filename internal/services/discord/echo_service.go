package discord_services

import (
	"fmt"
	"log"

	"bingobot/internal/models"
	utils "bingobot/internal/utils/discord"

	"github.com/bwmarrin/discordgo"
)

type DiscordEchoService struct {
	User *models.User
}

func (des DiscordEchoService) Handle(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	opts utils.OptionMap,
) {
	// Pass user data to service function
	message := opts["message"].StringValue()
	answer := fmt.Sprintf("**%s says:** %s", des.User.DiscordID, message)

	// Serialize the result and send via bot
	responseData := discordgo.InteractionResponseData{
		Content: answer,
	}
	response := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &responseData,
	}

	err := s.InteractionRespond(i.Interaction, &response)

	if err != nil {
		log.Panicf("could not respond to interaction: %s", err)
	}
}
