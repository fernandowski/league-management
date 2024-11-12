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
	ID             *TeamId
	Name           TeamName
	members        TeamMembers
	Roles          map[string]Role
	OrganizationId string
}

func NewTeamWithOwner(name TeamName, organizationId string, ownerId string) (*Team, error) {
	if strings.TrimSpace(string(name)) == "" {
		return nil, errors.New("name cannot be empty")
	}

	if strings.TrimSpace(string(organizationId)) == "" {
		return nil, errors.New("organization cannot be empty")
	}

	return &Team{ID: nil, Name: name, OrganizationId: organizationId, Roles: map[string]Role{ownerId: RoleOwner}}, nil
}

func NewTeamWithoutOwner(name TeamName, organizationId string) (*Team, error) {
	if strings.TrimSpace(string(name)) == "" {
		return nil, errors.New("name cannot be empty")
	}

	if strings.TrimSpace(organizationId) == "" {
		return nil, errors.New("organization cannot be empty")
	}

	return &Team{ID: nil, Name: name, OrganizationId: organizationId}, nil
}
