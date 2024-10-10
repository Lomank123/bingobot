package main

import (
	config "bingobot/configs/discord"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

type optionMap = map[string]*discordgo.ApplicationCommandInteractionDataOption

func main() {
	fmt.Println("Bingo Bot (Discord) starting...")
	config.LoadConfig()
	setupBot()
}

func parseOptions(options []*discordgo.ApplicationCommandInteractionDataOption) (om optionMap) {
	om = make(optionMap)

	for _, opt := range options {
		om[opt.Name] = opt
	}

	return
}

func parseAuthor(i *discordgo.Interaction) *discordgo.User {
	if i.Member != nil {
		return i.Member.User
	}

	return i.User
}

func handleEcho(s *discordgo.Session, i *discordgo.InteractionCreate, opts optionMap) {
	message := opts["message"].StringValue()
	author := parseAuthor(i.Interaction).String()
	answer := fmt.Sprintf("**%s says:** %s", author, message)

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

	log.Printf("Echo handled! %s", answer)
}

func setupBot() {
	session, _ := discordgo.New("Bot " + config.Cfg.DiscordBotToken)

	// Handlers
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			log.Printf("Invalid interaction type: %s", i.Type)
			return
		}

		data := i.ApplicationCommandData()

		if data.Name != "echo" {
			return
		}

		handleEcho(s, i, parseOptions(data.Options))
	})

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as %s", r.User.String())
	})

	// Commands
	_, err := session.ApplicationCommandBulkOverwrite(
		config.Cfg.DiscordAppId,
		"",
		config.Commands,
	)

	if err != nil {
		log.Fatalf("couldn't register commands: %s", err)
	}

	// Lifecycle?
	// TODO: Check what this does exactly
	err = session.Open()

	if err != nil {
		log.Fatalf("could not open session: %s", err)
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	err = session.Close()

	if err != nil {
		log.Printf("couldn't close session gracefully: %s", err)
	}
}
