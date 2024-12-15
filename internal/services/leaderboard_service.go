package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	consts "bingobot/internal/consts"
	"bingobot/internal/models"

	"github.com/redis/go-redis/v9"
)

// Leaderboard service based on Redis Sorted Set
type LeaderboardService struct {
	RedisClient *redis.Client
}

// Append user score to the recent leaderboard.
// If the leaderboard does not exist, it will be created.
func (ls LeaderboardService) RecordScore(user *models.User, score int) error {
	if score == 0 {
		return nil
	}

	dateStr := time.Now().Format(consts.LEADERBOARD_DATE_FORMAT)
	leaderboardID := fmt.Sprintf(consts.LEADERBOARD_ID_FORMAT, dateStr)

	_, err := ls.RedisClient.ZIncrBy(
		context.Background(),
		leaderboardID,
		float64(score),
		user.ID.String(),
	).Result()

	if err != nil {
		return err
	}

	return nil
}

// Return the leaderboard message for the given date range with user stats.
func (ls LeaderboardService) GetLeaderboardMessage(
	user *models.User,
	startDateStr, endDateStr string,
) (string, error) {
	err := ls.validateDates(startDateStr, endDateStr)

	if err != nil {
		return "Invalid dates format.", err
	}

	leaderboardID := ls.FormatLeaderboardID(startDateStr, endDateStr)
	var leaderboard *[]redis.Z

	// Here we will return leaderboard for the given month
	if startDateStr != "" && endDateStr == "" {
		leaderboard, err = ls.GetLeaderboardByID(leaderboardID)
	} else {
		leaderboard, err = ls.GetAggregatedLeaderboard(
			startDateStr, endDateStr, leaderboardID,
		)
	}

	if err != nil {
		log.Printf("could not get leaderboard: %s", err)
		return "Error occurred while getting the leaderboard. Try again later", err
	}

	userStats, err := ls.GetUserStats(user, leaderboardID)

	if err != nil {
		log.Printf("could not get user stats: %s", err)
		return "Error occurred while getting your stats. Try again later", err
	}

	return ls.buildLeaderboardMessage(
		leaderboard, userStats, startDateStr, endDateStr,
	), nil
}

// Return an aggregated leaderboard for the given dates range.
// If it doesn't exist, calculate and return it.
// Example (get all-time leaderboard):
// LeaderboardService.GetAggregatedLeaderboard("", "")
func (ls LeaderboardService) GetAggregatedLeaderboard(
	startDateStr, endDateStr string,
	aggregatedID string,
) (leaderboard *[]redis.Z, err error) {
	// TODO: Add limit param to limit the amount of entries returned from the leaderboard
	leaderboard, err = ls.GetLeaderboardByID(aggregatedID)

	if err != nil {
		return nil, err
	}

	currentDateStr := time.Now().Format(consts.LEADERBOARD_DATE_FORMAT)
	isCurrentMonth := currentDateStr == endDateStr

	// Return the leaderboard if it exists and end date is not up to present
	if len(*leaderboard) > 0 && !isCurrentMonth {
		return leaderboard, nil
	}

	leaderboardIDs, err := ls.getLeaderboardIDRange(startDateStr, endDateStr)

	if err != nil {
		return nil, err
	}

	err = ls.aggregateByIDRange(leaderboardIDs, aggregatedID)

	if err != nil {
		return nil, err
	}

	leaderboard, err = ls.GetLeaderboardByID(aggregatedID)

	if err != nil {
		return nil, err
	}

	return leaderboard, nil
}

// Return the user's rank and score in the given leaderboard.
func (ls LeaderboardService) GetUserStats(
	user *models.User,
	leaderboardID string,
) (*redis.RankScore, error) {
	rankWithScore, err := ls.RedisClient.ZRankWithScore(
		context.Background(),
		leaderboardID,
		user.ID.String(),
	).Result()

	if err != nil {
		return nil, err
	}

	return &rankWithScore, nil
}

func (ls LeaderboardService) GetLeaderboardByID(leaderboardID string) (
	*[]redis.Z, error,
) {
	leaderboard, err := ls.RedisClient.ZRevRangeWithScores(
		context.Background(),
		leaderboardID,
		0,
		-1,
	).Result()

	if err != nil {
		return nil, err
	}

	return &leaderboard, nil
}

// Construct the leaderboard ID based on the given date range.
// Possible formats: leaderboard:2021:01, leaderboard-aggregated-2021:01-2021:03
func (LeaderboardService) FormatLeaderboardID(
	startDateStr, endDateStr string,
) string {
	switch {
	// Return leaderboard ID for the given month
	case startDateStr != "" && endDateStr == "":
		return fmt.Sprintf(
			consts.LEADERBOARD_ID_FORMAT,
			startDateStr,
		)
	// Return all-time aggregated leaderboard ID
	case startDateStr == "" && endDateStr == "":
		return fmt.Sprintf(
			consts.LEADERBOARD_AGGREGATED_ID_FORMAT,
			consts.LEADERBOARD_START_DATE,
			time.Now().Format(consts.LEADERBOARD_DATE_FORMAT),
		)
	// Return aggregated leaderboard ID for the given date range
	default:
		return fmt.Sprintf(
			consts.LEADERBOARD_AGGREGATED_ID_FORMAT,
			startDateStr,
			endDateStr,
		)
	}
}

// Formats the leaderboard data for displaying in chat
func (LeaderboardService) buildLeaderboardMessage(
	leaderboard *[]redis.Z,
	userStats *redis.RankScore,
	startDateStr, endDateStr string,
) string {
	if len(*leaderboard) == 0 {
		return "Requested leaderboard has no records."
	}

	var sb strings.Builder
	var dateRange string

	switch {
	case startDateStr == "" && endDateStr == "":
		dateRange = "All time"
	case startDateStr != "" && endDateStr == "":
		dateRange = startDateStr
	default:
		dateRange = fmt.Sprintf("%s - %s", startDateStr, endDateStr)
	}

	sb.WriteString(fmt.Sprintf("Leaderboard (%s):\n\n", dateRange))

	for i, entry := range *leaderboard {
		sb.WriteString(fmt.Sprintf(
			"%d. %s - %d points\n",
			i+1,
			entry.Member,
			int(entry.Score),
		))
	}

	if userStats != nil {
		sb.WriteString(fmt.Sprintf(
			"\nYour rank in this leaderboard: %d\nYour score: %d points",
			userStats.Rank+1,
			int(userStats.Score),
		))
	}

	return sb.String()
}

// Aggregate scores for each user based on the leaderboard ID range.
// Writes it to the aggregateKey.
func (ls *LeaderboardService) aggregateByIDRange(
	leaderboardIDs []string,
	aggregateKey string,
) error {
	// Use ZUNIONSTORE to aggregate the sorted sets
	_, err := ls.RedisClient.ZUnionStore(
		context.Background(),
		aggregateKey,
		&redis.ZStore{
			Keys: leaderboardIDs,
		},
	).Result()

	if err != nil {
		return fmt.Errorf("failed to aggregate leaderboards: %w", err)
	}

	return nil
}

// Returns a sorted list of ids between the start and end dates.
// Example: ("2021:01", "2021:03") will return ["leaderboard:2021:01", "leaderboard:2021:02", "leaderboard:2021:03"], nil
func (LeaderboardService) getLeaderboardIDRange(
	startDateStr, endDateStr string,
) ([]string, error) {
	startDate, _ := time.Parse(consts.LEADERBOARD_DATE_FORMAT, startDateStr)
	endDate, _ := time.Parse(consts.LEADERBOARD_DATE_FORMAT, endDateStr)
	var dates []string

	// d.AddDate(0, 1, 0) adds 1 month to the date
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 1, 0) {
		leaderboardID := fmt.Sprintf(
			consts.LEADERBOARD_ID_FORMAT,
			d.Format(consts.LEADERBOARD_DATE_FORMAT),
		)
		dates = append(dates, leaderboardID)
	}

	return dates, nil
}

// Validate the given start and end dates.
func (LeaderboardService) validateDates(startDateStr, endDateStr string) error {
	// For all-time leaderboard
	if startDateStr == "" && endDateStr == "" {
		return nil
	}

	startDate, err := time.Parse(consts.LEADERBOARD_DATE_FORMAT, startDateStr)

	if err != nil {
		return fmt.Errorf("invalid start date format: %w", err)
	}

	// In case user wants to retrieve leaderboard for a single month
	if endDateStr == "" {
		return nil
	}

	// Here full date range is validated
	endDate, err := time.Parse(consts.LEADERBOARD_DATE_FORMAT, endDateStr)

	if err != nil {
		return fmt.Errorf("invalid end date format: %w", err)
	}

	// Ensure startDate is before or equal to endDate
	if startDate.After(endDate) {
		return fmt.Errorf("start date must be before or equal to end date")
	}

	return nil
}

// TODO: Add recalculation command (will be based on score service)

func NewLeaderboardService(redisClient *redis.Client) *LeaderboardService {
	return &LeaderboardService{RedisClient: redisClient}
}
