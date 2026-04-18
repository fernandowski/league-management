package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"league-management/internal/organization_management/domain/team"
	"league-management/internal/shared/database"
	"strings"
)

type TeamRepository struct {
}

func NewTeamRepository() *TeamRepository {
	return &TeamRepository{}
}

func (tr *TeamRepository) Save(team *team.Team) error {
	return database.WithTx(context.Background(), func(tx pgx.Tx) error {
		sql := "INSERT INTO league_management.teams (name) VALUES ($1) RETURNING id"

		var newTeamId string
		err := tx.QueryRow(context.Background(), sql, string(team.Name)).Scan(&newTeamId)
		if err != nil {
			return err
		}

		sql = "INSERT INTO league_management.organization_teams (organization_id, team_id) VALUES ($1, $2)"

		_, err = tx.Exec(context.Background(), sql, team.OrganizationId, newTeamId)
		if err != nil {
			return err
		}

		if len(team.Roles) > 0 {
			if err := upsertTeamOwnersTx(tx, newTeamId, team.Roles); err != nil {
				return err
			}
		}

		return nil
	})
}

func (tr *TeamRepository) FindById(teamId string) (*team.Team, error) {
	connection := database.GetConnection()

	sql := `
		SELECT
			teams.id as team_id,
			teams.name as team_name,
			organizations.id as organization_id,
			team_user_roles.user_id as user_id,
			team_user_roles.role as user_role
		FROM teams
		JOIN organization_teams ON organization_teams.team_id = teams.id
		JOIN organizations ON organizations.id = organization_teams.organization_id
		LEFT JOIN team_user_roles ON team_user_roles.team_id = teams.id
		WHERE
		teams.id=$1`

	rows, err := connection.Query(context.Background(), sql, teamId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var id, teamName, organizationId string
	var roles = make(map[string]team.TeamRole)

	if !rows.Next() {
		return nil, errors.New("team does not exist")
	}

	for rows.Next() {
		var roleUserId string
		var userRole string
		if err := rows.Scan(&id, &teamName, &organizationId, &roleUserId, &userRole); err != nil {
			return nil, err
		}
		roles[roleUserId] = team.TeamRole(userRole)
	}

	return &team.Team{
		ID:             (*team.TeamId)(&teamId),
		Name:           team.TeamName(teamName),
		OrganizationId: organizationId,
		Roles:          roles,
	}, nil
}

func (tr *TeamRepository) Search(organizationId, userId, searchTerm string) []interface{} {
	connection := database.GetConnection()
	queryBuilder := QueryBuilder{}

	sql := `SELECT teams.id, teams.name, organizations.name as organization_name 
		FROM organization_teams 
		JOIN organizations ON organizations.id = organization_teams.organization_id 
		JOIN teams on teams.id = organization_teams.team_id `

	queryBuilder.SetQuery(sql)
	queryBuilder.AddFilter("organization_teams.organization_id=$"+fmt.Sprint(len(queryBuilder.parameters)+1), organizationId)
	queryBuilder.AddFilter("organizations.user_id=$"+fmt.Sprint(len(queryBuilder.parameters)+1), userId)
	queryBuilder.AddFilter("teams.name ILIKE $"+fmt.Sprint(len(queryBuilder.parameters)+1), "%"+searchTerm+"%")
	//queryBuilder.SetPagination(1000, 0)
	query, parameters := queryBuilder.BuildQuery()

	rows, err := connection.Query(context.Background(), query, parameters...)

	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	results := []interface{}{}

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

	return results
}

func upsertTeamOwnersTx(tx pgx.Tx, teamId string, owners map[string]team.TeamRole) error {
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

	_, err := tx.Exec(context.Background(), sql)
	if err != nil {
		return err
	}

	return nil
}
