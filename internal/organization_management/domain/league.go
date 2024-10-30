package domain

type League struct {
	Id             string
	Name           string
	OwnerId        string
	OrganizationId string
	TeamIds        []string
}

func NewLeague(name string, ownerId string, organizationId string) League {
	return League{
		Name: name, OwnerId: ownerId, OrganizationId: organizationId, TeamIds: []string{},
	}
}
