package main

import (
	discord_config "bingobot/configs/discord"
	"log"
)

func main() {
	log.Println("Bingo Bot v1 starting...")
	cfg := discord_config.LoadConfig()
	log.Println(cfg.DBHost, cfg.DBName, cfg.DBUser, cfg.DBPassword)
}
