package repositories

import (
	"context"
	"fmt"
	"league-management/internal/organization_management/domain"
	"league-management/internal/shared/database"
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

	if len(team.Roles) > 0 {
		upsertTeamOwners(newTeamId, team.Roles)
	}

	return err
}

func upsertTeamOwners(teamId string, owners map[string]domain.Role) {
	var totalValues = len(owners)
	var values = make([]string, totalValues)

	index := 0
	for key, value := range owners {
		values[index] = fmt.Sprintf("('%s', '%s', '%s')", teamId, key, value)
		index++
	}

	sql := "INSERT INTO league_management.team_user_role (team_id, user_id, role) VALUES " +
		strings.Join(values, ",") +
		"ON CONFLICT DO NOTHING"

	connection := database.GetConnection()

	_, err := connection.Exec(context.Background(), sql)

	if err != nil {
		panic(err.Error())
	}
}
