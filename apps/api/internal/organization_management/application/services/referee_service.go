package services

import (
	"errors"
	seasonpkg "league-management/internal/organization_management/domain/season"
)

type RefereeService struct{}

func NewRefereeService() *RefereeService {
	return &RefereeService{}
}

// UpdateMatchScore allows a referee to update the score of a match through the season aggregate.
func (rs *RefereeService) UpdateMatchScore(currentSeason *seasonpkg.Season, matchID, refereeID string, homeScore, awayScore int) (*seasonpkg.Season, error) {
	if currentSeason == nil {
		return nil, errors.New("season cannot be nil")
	}
	if refereeID == "" {
		return nil, errors.New("refereeID cannot be empty")
	}
	if homeScore < 0 || awayScore < 0 {
		return nil, errors.New("scores must be non-negative")
	}

	return currentSeason.ChangeMatchScoreByReferee(matchID, refereeID, homeScore, awayScore)
}
