package discord_handlers

import (
	"log"

	consts "bingobot/internal/consts/discord"
	db "bingobot/internal/mongodb"
	services "bingobot/internal/services/discord"
	utils "bingobot/internal/utils/discord"

	"github.com/bwmarrin/discordgo"
)

func SetupHandlers(s *discordgo.Session) {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as %s", r.User.String())
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			log.Printf("Invalid interaction type: %s", i.Type)
			return
		}

		// TODO: Move DB collection access to service layer?
		// May help: https://medium.com/@shershnev/layered-architecture-implementation-in-golang-6318a72c1e10
		collection := db.DB.Collection(consts.USER_COLLECTION_NAME)
		userService := services.DiscordUserService{Collection: collection}
		user := userService.GetOrCreateUser(i.Interaction)
		data := i.ApplicationCommandData()
		options := utils.ParseOptions(data.Options)

		if data.Name == consts.ECHO_COMMAND {
			echoService := services.DiscordEchoService{User: user}
			echoService.Handle(s, i, options)
		}
	})
}
