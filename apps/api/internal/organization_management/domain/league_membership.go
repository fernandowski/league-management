package domain

import "errors"

type MembershipStatus string

const (
	MembershipActive    MembershipStatus = "Active"
	MembershipSuspended MembershipStatus = "Suspended"
	MembershipInactive  MembershipStatus = "Inactive"
)

type LeagueMembership struct {
	ID               string
	TeamID           string
	MemberShipStatus MembershipStatus
}

func NewLeagueMembership(MembershipId string, teamID string) (LeagueMembership, error) {
	return LeagueMembership{
		ID:               MembershipId,
		TeamID:           teamID,
		MemberShipStatus: MembershipInactive,
	}, nil
}

func (lm *LeagueMembership) Activate() (LeagueMembership, error) {

	if lm.MemberShipStatus == MembershipActive {
		return LeagueMembership{}, errors.New("membership already in active state")
	}

	return LeagueMembership{ID: lm.ID, TeamID: lm.TeamID, MemberShipStatus: MembershipActive}, nil
}
