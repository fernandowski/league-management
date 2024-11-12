package domainservices

import "league-management/internal/organization_management/domain"

func CreateTeam(teamName domain.TeamName, organization *domain.Organization, teamOwner string) (*domain.Team, error) {

	var newTeam *domain.Team
	var err error

	if organization.BelongsToOwner(teamOwner) {
		newTeam, err = domain.NewTeamWithoutOwner(teamName, *organization.ID)
	} else {
		newTeam, err = domain.NewTeamWithOwner(teamName, *organization.ID, teamOwner)
	}

	if err != nil {
		return nil, err
	}

	return newTeam, nil
}
