package discord_config

import (
	consts "bingobot/internal/consts/discord"

	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        consts.HELP_COMMAND,
		Description: "Shows tips for commands",
	},
	{
		Name:        consts.ECHO_COMMAND,
		Description: "Echoes your message back to the chat",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "message",
				Description: "Contents of the message",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	},
	{
		Name:        consts.MY_SCORE_COMMAND,
		Description: "Shows your current score",
	},
	{
		Name:        consts.LEADERBOARD_COMMAND,
		Description: "Shows the leaderboard based on the given date range. If no range then shows all-time leaderboard.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "start_date",
				Description: "Start date for the leaderboard. Format: YYYY-MM.",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
			{
				Name:        "end_date",
				Description: "End date for the leaderboard. Format: YYYY-MM.",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
		},
	},
}
