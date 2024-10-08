package discord_config

import (
	"log"

	utils "bingobot/internal"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost           string
	DBName           string
	DBUser           string
	DBPassword       string
	DiscordBotToken  string
	DiscordAppId     string
	DiscordPublicKey string
}

var Cfg Config

func LoadConfig() {
	// Can't be used in docker
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Cannot load .env")
	}

	Cfg.DBHost = utils.GetEnv("DB_HOST", "localhost")
	Cfg.DBName = utils.GetEnv("DB_NAME", "bingobot-db-1")
	Cfg.DBUser = utils.GetEnv("DB_USER", "bingobot-user-1")
	Cfg.DBPassword = utils.GetEnv("DB_PASSWORD", "bingobot-password-1")
	Cfg.DiscordBotToken = utils.GetEnv("DISCORD_BOT_TOKEN", "test-token-1")
	Cfg.DiscordAppId = utils.GetEnv("DISCORD_APP_ID", "app-id-1")
	Cfg.DiscordPublicKey = utils.GetEnv("DISCORD_PUBLIC_KEY", "public-key-1")
}