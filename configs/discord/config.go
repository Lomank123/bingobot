package discord_config

import (
	"fmt"
	"log"

	utils "bingobot/internal"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost           string
	DBPort           string
	DBName           string
	DBUser           string
	DBPassword       string
	DBURI            string
	DiscordBotToken  string
	DiscordAppId     string
	DiscordPublicKey string
	RedisHost        string
	RedisPort        string
	RedisAddress     string
	RedisUsername    string
	RedisPassword    string
}

var Cfg Config

func LoadConfig() {
	fmt.Println("Loading env...")

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Cannot load .env")
	}

	Cfg.DBHost = utils.GetEnv("DB_HOST", "localhost")
	Cfg.DBPort = utils.GetEnv("DB_PORT", "27017")
	Cfg.DBName = utils.GetEnv("DB_NAME", "bingobot-db-1")
	Cfg.DBUser = utils.GetEnv("DB_USER", "bingobot-user-1")
	Cfg.DBPassword = utils.GetEnv("DB_PASS", "bingobot-password-1")
	Cfg.DBURI = fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s",
		Cfg.DBUser,
		Cfg.DBPassword,
		Cfg.DBHost,
		Cfg.DBPort,
		Cfg.DBName,
	)
	Cfg.DiscordBotToken = utils.GetEnv("DISCORD_BOT_TOKEN", "test-token-1")
	Cfg.DiscordAppId = utils.GetEnv("DISCORD_APP_ID", "app-id-1")
	Cfg.DiscordPublicKey = utils.GetEnv("DISCORD_PUBLIC_KEY", "public-key-1")
	Cfg.RedisHost = utils.GetEnv("REDIS_HOST", "localhost")
	Cfg.RedisPort = utils.GetEnv("REDIS_PORT", "6379")
	Cfg.RedisAddress = fmt.Sprintf("%s:%s", Cfg.RedisHost, Cfg.RedisPort)
	Cfg.RedisUsername = utils.GetEnv("REDIS_USER", "redis-user-1")
	Cfg.RedisPassword = utils.GetEnv("REDIS_PASS", "redis-password-1")

	fmt.Println("Env loaded!")
}
