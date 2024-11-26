package domain

import (
	"errors"
	"strings"
)

type Season struct {
	ID       string
	LeagueId string
	Name     string
}

func NewSeason(name, leagueID string) (*Season, error) {
	if strings.TrimSpace(leagueID) == "" {
		return nil, errors.New("leagueId cannot be empty")
	}

	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name cannot be empty")
	}

	return &Season{ID: "", Name: name, LeagueId: leagueID}, nil
}
