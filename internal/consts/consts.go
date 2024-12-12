package consts

// Domains

const TELEGRAM_DOMAIN = "telegram"
const DISCORD_DOMAIN = "discord"

// Collections

const USER_COLLECTION_NAME = "users"
const USER_SCORE_RECORD_COLLECTION_NAME = "user_score_records"
const USER_TELEGRAM_PROFILE_COLLECTION_NAME = "user_telegram_profiles"
const USER_DISCORD_PROFILE_COLLECTION_NAME = "user_discord_profiles"

// Response texts

var COMMAND_NOT_FOUND_TEXT = "I don't know that command. Try /help to see all available commands."

// Score

var SCORE_LIMIT_PER_DAY = 100

// Leaderboard

// Date of the oldest possible leaderboard
var LEADERBOARD_START_DATE = "2000:01"
// Date format: "YYYY:MM"
var LEADERBOARD_DATE_FORMAT = "2006:01"
// Example: "leaderboard:2021:01"
var LEADERBOARD_ID_FORMAT = "leaderboard:%s"
// Example: "leaderboard-aggregated-2021:01-2021:03"
var LEADERBOARD_AGGREGATED_ID_FORMAT = "leaderboard-aggregated-%s-%s"
