package services

import (
	"league-management/internal/organization_management/domain/organization"
	"league-management/internal/organization_management/domain/team"
)

func CreateTeam(teamName team.TeamName, organization *organization.Organization, teamOwner string) (*team.Team, error) {

	var newTeam *team.Team
	var err error

	if organization.BelongsToOwner(teamOwner) {
		newTeam, err = team.NewTeamWithoutOwner(teamName, *organization.ID)
	} else {
		newTeam, err = team.NewTeamWithOwner(teamName, *organization.ID, teamOwner)
	}

	if err != nil {
		return nil, err
	}

	return newTeam, nil
}
