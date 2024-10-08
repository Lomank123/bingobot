# bingobot

TG/Discord bot with various unique features

## Prerequisites

- Go (^1.22.4)
- Docker, docker compose
- k8s


## Telegram

- Open Telegram
- Find user named "BotFather"
- Start a chat (`/start`)
- Follow the instructions and in the end copy the generated API key to `TELEGRAM_BOT_TOKEN` variable in `.env`


## Discord

- Go to https://discord.com/developers/applications
- Create your application
- Go to "Bot" section
- Click on "Reset Token" button
- Copy generated bot token to .env (`DISCORD_BOT_TOKEN` variable)
- Go to "General Information" section
- Copy "Application ID" and "Public key" to corresponding env variables