# bingobot

TG/Discord bot with various unique features


## Roadmap

### Main features

- Everything is shared between discord and tg bots

Actions:
- Scoring system
  - User receives points after each command
  - Store the score from TG and Discord separately
- Real-time Score Leaderboard (based on scoring system)
  - Think about system design in advance
- Bets (based on scoring system)
  - Admin writes a question with options to answer
  - User can deposit up to a certain amount of points
- Achievements
  - Think of system design
  - Use scoring system
  - Use any other type of activity


### Future plans

- Scoring system
  - Add daily limit
  - Add points transfer with limitations
- Unique TG game (triggered by 1+ users)
  - Game Statistics
- Unique Discord game (triggered by 1+ users)
  - Game Statistics
- Keep conversation and save some of the quotes. Send these quotes later after ~100 messages in chat
- Guess who's message game (send random saved message/quote and set timer. Collect messages after it and if there's a correct nickname of the user then it guessed correctly. Otherwise after timer has expired tell the right answer)
- Paid subscription


### DND ideas

- d2, d4, d6 flips
- Save biography of a character (can be used later)
- Save notes about the campaign
- Spellbook (requires API of spells and editions of dnd)


## Prerequisites

- Go (^1.22.4)
- Docker & docker compose 
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


## Installation

- Clone the repo:

```shell
git clone https://github.com/Lomank123/bingobot.git
```

- Copy `.env.sample` to `.env` and fill the required env variables:

```shell
cp .env.sample .env
```

- Create docker network:

```shell
docker network create bingo-network
```

- Build and run the containers:

```shell
docker compose up -d --build
```

- In case you need production build:

```shell
docker compose -f docker-compose-prod.yml up -d --build
```
