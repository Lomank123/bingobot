package discord_handlers

import (
	"fmt"
	"log"
	"time"

	general_consts "bingobot/internal/consts"
	consts "bingobot/internal/consts/discord"
	"bingobot/internal/models"
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

		user, err := getOrCreateUser(srvs, i.Interaction)

		if err != nil {
			log.Panicf("could not get or create user: %s", err)
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
			score, err := srvs.ScoreService.GetUserTotalScore(user)

			if err != nil {
				log.Printf("could not get user score: %s", err)
				message = "Error occurred while getting your score. Try again later"
			}

			message = fmt.Sprintf("Your total score is %d points. Well done!", score)
		case consts.LEADERBOARD_COMMAND:
			var startDateStr string
			var endDateStr string

			if val, ok := options["start_date"]; ok {
				startDateStr = val.StringValue()
			}
			if val, ok := options["end_date"]; ok {
				endDateStr = val.StringValue()
			}

			message, _ = srvs.LeaderboardService.GetLeaderboardMessage(
				user,
				startDateStr,
				endDateStr,
				general_consts.DISCORD_DOMAIN,
			)
		default:
			message = general_consts.COMMAND_NOT_FOUND_TEXT
		}

		// Record the score
		score, err := srvs.ScoreService.RecordScore(
			user,
			data.Name,
			general_consts.DISCORD_DOMAIN,
		)

		if err != nil {
			log.Printf("could not record score: %s", err)
		}

		err = srvs.LeaderboardService.RecordScore(
			user.ID.Hex(), score, time.Now(),
		)

		if err != nil {
			// TODO: Perhaps we need to check how much this error occurred.
			// If more than 3 times, then we have to trigger re-calculation.
			log.Printf("could not record score for leaderboard: %s", err)
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

// TODO: Think of moving such functionality to a separate service
func getOrCreateUser(
	srvs *services.DiscordService,
	i *discordgo.Interaction,
) (*models.User, error) {
	// TODO: Make inserts in a single transaction
	discordUser := utils.ParseDiscordUser(i)
	user, isCreated, err := srvs.UserService.GetOrCreate(discordUser.ID, "")

	if isCreated {
		_, err = srvs.ProfileService.Create(user, discordUser)

		if err != nil {
			log.Printf("could not create user profile: %s", err)
			return nil, err
		}
	}

	if err != nil {
		log.Printf("error occurred during user retrieval/creation: %s", err)
		return nil, err
	}

	return user, nil
}
