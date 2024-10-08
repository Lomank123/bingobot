package telegram_config

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
	TelegramBotToken string
}

var Cfg Config

func LoadConfig() {
	// TODO: Can't be used in docker
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Cannot load .env")
	}

	Cfg.DBHost = utils.GetEnv("DB_HOST", "localhost")
	Cfg.DBName = utils.GetEnv("DB_NAME", "bingobot-db-1")
	Cfg.DBUser = utils.GetEnv("DB_USER", "bingobot-user-1")
	Cfg.DBPassword = utils.GetEnv("DB_PASSWORD", "bingobot-password-1")
	Cfg.TelegramBotToken = utils.GetEnv("TELEGRAM_BOT_TOKEN", "test-token-1")
}
