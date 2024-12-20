package leaderboard_utils

import (
	services "bingobot/internal/services"
	"fmt"
)

// Aggregate leaderboard data from the DB, clear all previous leaderboards from Redis,
// calculate the leaderboards and store them in Redis Sorted Sets.
func ResetLeaderboards(
	scoreService *services.ScoreService,
	leaderboardService *services.LeaderboardService,
) error {
	fmt.Printf("Aggregating leaderboard data from DB...")
	scores, err := scoreService.AggregateLeaderboardData()

	if err != nil {
		return fmt.Errorf("error aggregating leaderboard data: %s", err)
	}

	fmt.Printf("Clearing all previous leaderboards from Redis...")
	err = leaderboardService.ClearAllLeaderboards()

	if err != nil {
		return fmt.Errorf("error clearing the leaderboards: %s", err)
	}

	fmt.Printf("Start calculating the leaderboards...")
	err = leaderboardService.RecalculateLeaderboard(scores)

	if err != nil {
		return fmt.Errorf("error calculating the leaderboards: %s", err)
	}

	fmt.Printf("Leaderboards have been calculated successfully!")

	return nil
}
