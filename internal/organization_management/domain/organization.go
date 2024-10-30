package domain

import "errors"

type Organization struct {
	ID                  *string
	Name                string
	slug                string
	OrganizationOwnerId string
	isActive            bool
	DateCreated         string
	DateUpdated         string
}

func NewOrganization(id *string, name string, ownerId string, isActive bool) Organization {
	return Organization{ID: id, Name: name, OrganizationOwnerId: ownerId, isActive: isActive}
}

func (o *Organization) isInGoodStanding() bool {
	return o.isActive
}

func (o *Organization) BelongsToOwner(ownerId string) bool {
	return o.OrganizationOwnerId == ownerId
}

func (o *Organization) CanAcceptANewLeague(ownerId string) error {

	if !o.isInGoodStanding() {
		return errors.New("league not in good standing")
	}

	if !o.BelongsToOwner(ownerId) {
		return errors.New("not allowed to add league")
	}

	return nil
}
