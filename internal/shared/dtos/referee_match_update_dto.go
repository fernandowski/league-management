package dtos

type RefereeMatchUpdateDTO struct {
	MatchID   string `json:"match_id"`
	RefereeID string `json:"referee_id"`
	HomeScore int    `json:"home_score"`
	AwayScore int    `json:"away_score"`
}
