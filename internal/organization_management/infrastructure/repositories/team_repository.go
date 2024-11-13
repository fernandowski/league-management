package repositories

import (
	"context"
	"fmt"
	"league-management/internal/organization_management/domain"
	"league-management/internal/shared/database"
	"log"
	"strings"
)

type TeamRepository struct {
}

func NewTeamRepository() *TeamRepository {
	return &TeamRepository{}
}

func (tr *TeamRepository) Save(team *domain.Team) error {
	connection := database.GetConnection()

	sql := "INSERT INTO league_management.teams (name) VALUES ($1) RETURNING id"

	var newTeamId string
	err := connection.QueryRow(context.Background(), sql, string(team.Name)).Scan(&newTeamId)

	if err != nil {
		panic(err.Error())
	}

	sql = "INSERT INTO league_management.organization_teams (organization_id, team_id) VALUES ($1, $2)"

	_, err = connection.Exec(context.Background(), sql, team.OrganizationId, newTeamId)

	if err != nil {
		panic(err.Error())
	}

	if len(team.Roles) > 0 {
		upsertTeamOwners(newTeamId, team.Roles)
	}

	return err
}

func (tr *TeamRepository) FetchAll(organizationId string, userId string) []interface{} {
	connection := database.GetConnection()

	sql := "SELECT teams.id, teams.name, organizations.name as organization_name " +
		"FROM organization_teams " +
		"JOIN organizations ON organizations.id = organization_teams.organization_id " +
		"JOIN teams on teams.id = organization_teams.team_id " +
		"WHERE organization_teams.organization_id=$1 AND organizations.user_id= $2 "

	rows, err := connection.Query(context.Background(), sql, organizationId, userId)

	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	var results []interface{}

	for rows.Next() {
		var teamId, teamName, organizationName string
		if err := rows.Scan(&teamId, &teamName, &organizationName); err != nil {
			return results
		}

		team := make(map[string]string)

		team["name"] = teamName
		team["id"] = teamId
		team["organization_name"] = organizationName
		results = append(results, team)
	}

	log.Print("res", results)
	return results
}

func upsertTeamOwners(teamId string, owners map[string]domain.Role) {
	var totalValues = len(owners)
	var values = make([]string, totalValues)

	index := 0
	for key, value := range owners {
		values[index] = fmt.Sprintf("('%s', '%s', '%s')", teamId, key, value)
		index++
	}

	sql := "INSERT INTO league_management.team_user_roles (team_id, user_id, role) VALUES " +
		strings.Join(values, ",") +
		"ON CONFLICT (team_id, user_id, role) DO NOTHING"

	connection := database.GetConnection()

	_, err := connection.Exec(context.Background(), sql)

	if err != nil {
		panic(err.Error())
	}
}
