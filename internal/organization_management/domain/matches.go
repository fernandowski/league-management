package domain

import (
	"errors"
	"github.com/google/uuid"
)

type Match struct {
	ID               string
	HomeTeamID       string
	AwayTeamID       string
	HomeTeamScore    int
	AwayTeamScore    int
	AssignedLocation MatchLocation
}

func NewMatch(homeTeamID, awayTeamID string) (Match, error) {
	if homeTeamID == "" {
		return Match{}, errors.New("valid match must have a valid home team ID")
	}

	if awayTeamID == "" {
		return Match{}, errors.New("valid match must have a valid away team ID")
	}

	return Match{
		ID:            uuid.New().String(),
		HomeTeamID:    homeTeamID,
		AwayTeamID:    awayTeamID,
		HomeTeamScore: 0,
		AwayTeamScore: 0,
	}, nil
}

func (m *Match) GetHomeTeam() interface{} {
	if m.HomeTeamID == "bye" {
		return nil
	}
	return m.HomeTeamID
}

func (m *Match) GetAwayTeam() interface{} {
	if m.AwayTeamID == "bye" {
		return nil
	}
	return m.AwayTeamID
}
