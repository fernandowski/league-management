package domain

import (
	"errors"
	"github.com/google/uuid"
	"log"
)

type Match struct {
	ID               string
	HomeTeamID       string
	AwayTeamID       string
	HomeTeamScore    int
	AwayTeamScore    int
	AssignedLocation MatchLocation
}

func NewMatch(matchID *string, homeTeamID, awayTeamID string) (Match, error) {
	if homeTeamID == "" {
		return Match{}, errors.New("valid match must have a valid home team ID")
	}

	if awayTeamID == "" {
		return Match{}, errors.New("valid match must have a valid away team ID")
	}

	id := uuid.New().String()
	if matchID != nil {
		id = *matchID
	}

	return Match{
		ID:            id,
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

func (m *Match) ChangeScore(homeTeamScore, awayTeamScore int) (*Match, error) {
	if homeTeamScore < 0 {
		return nil, errors.New("score must be greater than 0")
	}

	if awayTeamScore < 0 {
		return nil, errors.New("score must be greater than 0")
	}

	log.Print(m)
	if m.HomeTeamID == "bye" || m.AwayTeamID == "bye" {
		return nil, errors.New("cannot change score of match for a bye week")
	}

	newMatch := *m
	newMatch.AwayTeamScore = awayTeamScore
	newMatch.HomeTeamScore = homeTeamScore

	return &newMatch, nil
}
