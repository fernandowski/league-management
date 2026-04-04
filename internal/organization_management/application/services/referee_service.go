package services

import (
	"errors"
	"league-management/internal/organization_management/domain"
)

type RefereeService struct{}

func NewRefereeService() *RefereeService {
	return &RefereeService{}
}

// UpdateMatchScore allows a referee to update the score of a match
func (rs *RefereeService) UpdateMatchScore(match *domain.Match, refereeID string, homeScore, awayScore int) error {
	if match == nil {
		return errors.New("match cannot be nil")
	}
	if refereeID == "" {
		return errors.New("refereeID cannot be empty")
	}
	if homeScore < 0 || awayScore < 0 {
		return errors.New("scores must be non-negative")
	}
	match.HomeTeamScore = homeScore
	match.AwayTeamScore = awayScore
	match.RefereeID = refereeID
	return nil
}
