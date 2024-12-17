package types

// Represent aggregated sum of user score for a given year and month
type AggregatedUserScore struct {
	UserID string `bson:"user_id"`
	Year   int    `bson:"year"`
	Month  int    `bson:"month"`
	Score  int    `bson:"score"`
}
