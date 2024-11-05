package domain

import (
	"errors"
	"strings"
)

type TeamId string
type TeamMembers []string
type TeamName string

type Role string

const (
	RoleOwner          Role = "owner"
	RoleCoach          Role = "coach"
	RoleAssistantCoach Role = "assistant_coach"
	RolePlayer         Role = "Player"
)

type Team struct {
	ID      *TeamId
	Name    TeamName
	members TeamMembers
	Roles   map[string]Role
}

func NewTeam(name TeamName, ownerId string) (*Team, error) {
	if strings.TrimSpace(string(name)) == "" {
		return nil, errors.New("name cannot be empty")
	}

	return &Team{ID: nil, Name: name, Roles: map[string]Role{ownerId: RoleOwner}}, nil
}
