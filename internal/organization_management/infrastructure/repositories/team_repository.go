package repositories

import (
	"context"
	"errors"
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

func (tr *TeamRepository) FindById(teamId string) (*domain.Team, error) {
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
	var roles = make(map[string]domain.TeamRole)

	if !rows.Next() {
		return nil, errors.New("team does not exist")
	}

	for rows.Next() {
		var roleUserId string
		var userRole string
		if err := rows.Scan(&id, &teamName, &organizationId, &roleUserId, &userRole); err != nil {
			return nil, err
		}
		roles[roleUserId] = domain.TeamRole(userRole)
	}

	return &domain.Team{
		ID:             (*domain.TeamId)(&teamId),
		Name:           domain.TeamName(teamName),
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

	log.Print(query)
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

func upsertTeamOwners(teamId string, owners map[string]domain.TeamRole) {
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
