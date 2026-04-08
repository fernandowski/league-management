package domain

import (
	"errors"
	"github.com/google/uuid"
)

type League struct {
	Id             *string
	Name           string
	OwnerId        string
	OrganizationId string
	ActiveSeason   string
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

func (l *League) RemoveMembership(membershipId string) (*League, error) {
	if len(l.Memberships) <= 0 {
		return nil, errors.New("membership empty")
	}

	league := l
	memberships := []LeagueMembership{}

	for _, membership := range l.Memberships {
		if membership.ID != membershipId {
			memberships = append(memberships, membership)
		}
	}

	league.Memberships = memberships

	return league, nil
}

func (l *League) HasActiveSeason() bool {
	if l.ActiveSeason == "" {
		return false
	}
	return true
}
