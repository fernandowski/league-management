package domain

import (
	"errors"
	"strings"
)

type TeamId string
type TeamMembers []string
type LeagueName string

type Team struct {
	ID      *TeamId
	Name    LeagueName
	members TeamMembers
}

func NewTeam(name LeagueName) (*Team, error) {
	if strings.TrimSpace(string(name)) == "" {
		return nil, errors.New("name cannot be empty")
	}

	return &Team{ID: nil, Name: name}, nil
}
