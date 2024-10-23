package domain

type Organization struct {
	ID                  string
	name                string
	slug                string
	OrganizationOwnerId string
}

func NewOrganization(id string, name string, ownerId string) Organization {
	return Organization{ID: id, name: name, OrganizationOwnerId: ownerId}
}
