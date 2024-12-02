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
}
