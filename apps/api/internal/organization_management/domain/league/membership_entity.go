package league

import "errors"

type MembershipStatus string

const (
	MembershipActive    MembershipStatus = "Active"
	MembershipSuspended MembershipStatus = "Suspended"
	MembershipInactive  MembershipStatus = "Inactive"
)

type Membership struct {
	ID               string
	TeamID           string
	MemberShipStatus MembershipStatus
}

func NewMembership(MembershipId string, teamID string) (Membership, error) {
	return Membership{
		ID:               MembershipId,
		TeamID:           teamID,
		MemberShipStatus: MembershipInactive,
	}, nil
}

func (lm *Membership) Activate() (Membership, error) {

	if lm.MemberShipStatus == MembershipActive {
		return Membership{}, errors.New("membership already in active state")
	}

	return Membership{ID: lm.ID, TeamID: lm.TeamID, MemberShipStatus: MembershipActive}, nil
}
