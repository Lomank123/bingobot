package discord_handlers

import (
	"log"

	general_consts "bingobot/internal/consts"
	consts "bingobot/internal/consts/discord"
	services "bingobot/internal/services/discord"
	utils "bingobot/internal/utils/discord"

	"github.com/bwmarrin/discordgo"
)

func SetupHandlers(s *discordgo.Session, srvs *services.DiscordService) {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as %s", r.User.String())
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			log.Printf("Invalid interaction type: %s", i.Type)
			return
		}

		user, _, err := srvs.UserService.GetOrCreateUser(i.Interaction)

		if err != nil {
			log.Panicf("error occurred during user retrieval/creation: %s", err)
		}

		data := i.ApplicationCommandData()
		options := utils.ParseOptions(data.Options)
		message := ""

		switch data.Name {
		case consts.ECHO_COMMAND:
			message = srvs.EchoService.Handle(options, user)
		case consts.HELP_COMMAND:
			// TODO: Implement command
			message = general_consts.COMMAND_NOT_FOUND_TEXT
		case consts.MY_SCORE_COMMAND:
			// TODO: Implement command
			message = general_consts.COMMAND_NOT_FOUND_TEXT
		case consts.LEADERBOARD_COMMAND:
			// TODO: Implement command
			message = general_consts.COMMAND_NOT_FOUND_TEXT
		default:
			message = general_consts.COMMAND_NOT_FOUND_TEXT
		}

		err = srvs.ScoreService.IncrementScore(user, data.Name)

		if err != nil {
			log.Printf("could not increment score: %s", err)
		}

		// Serialize the result and send via bot
		responseData := discordgo.InteractionResponseData{
			Content: message,
		}
		response := discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &responseData,
		}
		err = s.InteractionRespond(i.Interaction, &response)

		if err != nil {
			log.Panicf("could not respond to interaction: %s", err)
		}
	})
}
