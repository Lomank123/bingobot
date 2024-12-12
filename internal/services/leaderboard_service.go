package services

import (
	"context"
	"fmt"
	"time"

	consts "bingobot/internal/consts"

	"github.com/redis/go-redis/v9"
)

// Leaderboard service based on Redis Sorted Set
type LeaderboardService struct {
	RedisClient *redis.Client
}

// Append user score to the recent leaderboard.
// If the leaderboard does not exist, it will be created.
func (ls LeaderboardService) RecordScore(userID string, score int) error {
	dateStr := time.Now().Format(consts.LEADERBOARD_DATE_FORMAT)
	leaderboardID := fmt.Sprintf(consts.LEADERBOARD_ID_FORMAT, dateStr)

	_, err := ls.RedisClient.ZIncrBy(
		context.Background(),
		leaderboardID,
		float64(score),
		userID,
	).Result()

	if err != nil {
		return err
	}

	return nil
}

func (ls LeaderboardService) GetLeaderboardByID(leaderboardID string) ([]redis.Z, error) {
	leaderboard, err := ls.RedisClient.ZRevRangeWithScores(
		context.Background(),
		leaderboardID,
		0,
		-1,
	).Result()

	if err != nil {
		return nil, err
	}

	return leaderboard, nil
}

// Return an aggregated leaderboard for the given dates range.
// If isAllTime is true, then it will return the all-time leaderboard.
// If it doesn't exist, calculate and return it.
func (ls LeaderboardService) GetAggregatedLeaderboard(
	startDateStr, endDateStr string,
	isAllTime bool,
) ([]redis.Z, error) {
	aggregateKey := fmt.Sprintf(
		consts.LEADERBOARD_AGGREGATED_ID_FORMAT,
		startDateStr,
		endDateStr,
	)
	currentDateStr := time.Now().Format(consts.LEADERBOARD_DATE_FORMAT)

	if isAllTime {
		aggregateKey = fmt.Sprintf(
			consts.LEADERBOARD_AGGREGATED_ID_FORMAT,
			consts.LEADERBOARD_START_DATE,
			currentDateStr,
		)
	}

	leaderboard, err := ls.GetLeaderboardByID(aggregateKey)

	if err != nil {
		return nil, err
	}

	isCurrentMonth := currentDateStr == endDateStr

	// Return the leaderboard if it exists and end date is not up to present
	if len(leaderboard) > 0 && !isCurrentMonth {
		return leaderboard, nil
	}

	leaderboardIDs, err := ls.getLeaderboardIDRange(startDateStr, endDateStr)

	if err != nil {
		return nil, err
	}

	aggregateKey, err = ls.AggregateByIDRange(leaderboardIDs, aggregateKey)

	if err != nil {
		return nil, err
	}

	leaderboard, err = ls.GetLeaderboardByID(aggregateKey)

	if err != nil {
		return nil, err
	}

	return leaderboard, nil
}

// Aggregate scores for each user based on the leaderboard ID range.
// Returns new leaderboard ID.
func (ls *LeaderboardService) AggregateByIDRange(
	leaderboardIDs []string,
	aggregateKey string,
) (string, error) {
	// Use ZUNIONSTORE to aggregate the sorted sets
	_, err := ls.RedisClient.ZUnionStore(
		context.Background(),
		aggregateKey,
		&redis.ZStore{
			Keys: leaderboardIDs,
		},
	).Result()

	if err != nil {
		return "", fmt.Errorf("failed to aggregate leaderboards: %w", err)
	}

	return aggregateKey, nil
}

// Returns a sorted list of ids between the start and end dates.
// Example: ("2021:01", "2021:03") will return ["leaderboard:2021:01", "leaderboard:2021:02", "leaderboard:2021:03"], nil
func (LeaderboardService) getLeaderboardIDRange(
	startDateStr,
	endDateStr string,
) ([]string, error) {
	// Parse the start and end dates
	startDate, err := time.Parse(consts.LEADERBOARD_DATE_FORMAT, startDateStr)

	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	endDate, err := time.Parse(consts.LEADERBOARD_DATE_FORMAT, endDateStr)

	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	// Ensure startDate is before or equal to endDate
	if startDate.After(endDate) {
		return nil, fmt.Errorf("start date must be before or equal to end date")
	}

	// Generate the list of dates
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

// TODO: Add GetLeaderboardByUser function to return user's rank and score
// TODO: Add recalculation command (will be based on score service)

func NewLeaderboardService(redisClient *redis.Client) *LeaderboardService {
	return &LeaderboardService{RedisClient: redisClient}
}
