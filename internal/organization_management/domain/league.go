package domain

type League struct {
	Id             *string
	Name           string
	OwnerId        string
	OrganizationId string
	TeamIds        []string
}

func NewLeague(id string, name string, ownerId string, organizationId string) League {
	return League{
		Id: &id, Name: name, OwnerId: ownerId, OrganizationId: organizationId, TeamIds: []string{},
	}
}
