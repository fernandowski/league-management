package dtos

type ChangeGameScoreDTO struct {
	MatchID   string `json:"match_id"`
	HomeScore int    `json:"home_score"`
	AwayScore int    `json:"away_score"`
}
