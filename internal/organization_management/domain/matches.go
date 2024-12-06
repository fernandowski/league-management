package domain

import "errors"

type Match struct {
	HomeTeamID       string
	AwayTeamID       string
	AssignedLocation MatchLocation
}

func NewMatch(homeTeamID, awayTeamID string) (Match, error) {
	if homeTeamID == "" {
		return Match{}, errors.New("valid match must have a valid home team ID")
	}

	if awayTeamID == "" {
		return Match{}, errors.New("valid match must have a valid away team ID")
	}

	return Match{HomeTeamID: homeTeamID, AwayTeamID: awayTeamID}, nil
}
