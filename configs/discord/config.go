package discord_config

import (
	"log"
	"os"

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

	Cfg.DBHost = getEnv("DB_HOST", "localhost")
	Cfg.DBName = getEnv("DB_NAME", "bingobot-db-1")
	Cfg.DBUser = getEnv("DB_USER", "bingobot-user-1")
	Cfg.DBPassword = getEnv("DB_PASSWORD", "bingobot-password-1")
	Cfg.DiscordBotToken = getEnv("DISCORD_BOT_TOKEN", "test-token-1")
	Cfg.DiscordAppId = getEnv("DISCORD_APP_ID", "app-id-1")
	Cfg.DiscordPublicKey = getEnv("DISCORD_PUBLIC_KEY", "public-key-1")
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	log.Printf("WARNING: Default value '%s' is used for key: '%s'", defaultVal, key)
	return defaultVal
}
