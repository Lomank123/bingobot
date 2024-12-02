package discord_consts

// Commands

var ECHO_COMMAND = "echo"
var HELP_COMMAND = "help"
var MY_SCORE_COMMAND = "score"
var LEADERBOARD_COMMAND = "leaderboard"

var COMMAND_SCORE_MAPPING = map[string]int{
	ECHO_COMMAND:        1,
	HELP_COMMAND:        0,
	MY_SCORE_COMMAND:    0,
	LEADERBOARD_COMMAND: 0,
}
