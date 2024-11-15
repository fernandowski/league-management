package domain

import (
	"github.com/google/uuid"
)

type League struct {
	Id             *string
	Name           string
	OwnerId        string
	OrganizationId string
	Memberships    []LeagueMembership
}

func NewLeague(id string, name string, ownerId string, organizationId string) League {
	return League{
		Id: &id, Name: name, OwnerId: ownerId, OrganizationId: organizationId, Memberships: []LeagueMembership{},
	}
}

func (l *League) StartTeamMembership(teamId string) (*LeagueMembership, error) {
	newLeagueMembership, err := NewLeagueMembership(uuid.New().String(), teamId)

	if _, err := newLeagueMembership.Activate(); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	l.Memberships = append(l.Memberships, newLeagueMembership)

	return &newLeagueMembership, nil
}
