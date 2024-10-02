package discord_config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBName     string
	DBUser     string
	DBPassword string
}

func LoadConfig() *Config {
	// Can't be used in docker
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Cannot load .env")
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBName:     getEnv("DB_NAME", "bingobot-db-1"),
		DBUser:     getEnv("DB_USER", "bingobot-user-1"),
		DBPassword: getEnv("DB_PASSWORD", "bingobot-password-1"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	log.Printf("WARNING: Default value '%s' is used for key: '%s'", defaultVal, key)
	return defaultVal
}
