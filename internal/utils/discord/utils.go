package discord_utils

import "github.com/bwmarrin/discordgo"

type OptionMap = map[string]*discordgo.ApplicationCommandInteractionDataOption

func ParseOptions(options []*discordgo.ApplicationCommandInteractionDataOption) (om OptionMap) {
	om = make(OptionMap)

	for _, opt := range options {
		om[opt.Name] = opt
	}

	return
}

func ParseDiscordUser(i *discordgo.Interaction) *discordgo.User {
	if i.Member != nil {
		return i.Member.User
	}

	return i.User
}
