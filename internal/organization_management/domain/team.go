package domain

import (
	"errors"
	"strings"
)

type TeamId string
type TeamMembers []string
type TeamName string

type TeamRole string

const (
	RoleOwner          TeamRole = "owner"
	RoleCoach          TeamRole = "coach"
	RoleAssistantCoach TeamRole = "assistant_coach"
	RolePlayer         TeamRole = "Player"
)

type Team struct {
	ID             *TeamId
	Name           TeamName
	members        TeamMembers
	Roles          map[string]TeamRole
	OrganizationId string
}

func NewTeamWithOwner(name TeamName, organizationId string, ownerId string) (*Team, error) {
	if strings.TrimSpace(string(name)) == "" {
		return nil, errors.New("name cannot be empty")
	}

	if strings.TrimSpace(string(organizationId)) == "" {
		return nil, errors.New("organization cannot be empty")
	}

	return &Team{ID: nil, Name: name, OrganizationId: organizationId, Roles: map[string]TeamRole{ownerId: RoleOwner}}, nil
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
