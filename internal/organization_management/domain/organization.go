package domain

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

func (o *Organization) IsInGoodStanding() bool {
	return o.isActive
}
