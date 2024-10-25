package domain

type Organization struct {
	ID                  string
	Name                string
	slug                string
	OrganizationOwnerId string
	DateCreated         string
	DateUpdated         string
}

func NewOrganization(name string, ownerId string) Organization {
	return Organization{Name: name, OrganizationOwnerId: ownerId}
}
