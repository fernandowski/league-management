package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"league-management/internal/organization_management/domain"
	"league-management/internal/shared/database"
	"strings"
	"time"
)

type LeagueRepository struct {
}

func NewLeagueRepository() LeagueRepository {
	return LeagueRepository{}
}

func (lr *LeagueRepository) FindById(leagueId string) (*domain.League, error) {
	connection := database.GetConnection()

	sql := "SELECT id, name, user_id, organization_id, created_at, updated_at " +
		"FROM league_management.leagues " +
		"WHERE id=$1"

	var id string
	var name string
	var ownerId string
	var organizationId string
	var dateCreated time.Time
	var dateUpdated time.Time

	err := connection.QueryRow(context.Background(), sql, leagueId).Scan(&id, &name, &ownerId, &organizationId, &dateCreated, &dateUpdated)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		panic(err.Error())
	}

	sql = `SELECT id, team_id FROM league_teams WHERE league_id=$1`

	rows, err := connection.Query(context.Background(), sql, leagueId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	memberships := []domain.LeagueMembership{}

	for rows.Next() {
		var id string
		var teamId string
		if err := rows.Scan(&id, &teamId); err != nil {

		}
		memberships = append(memberships, domain.LeagueMembership{ID: id, TeamID: teamId})
	}

	return &domain.League{
		Id:             &id,
		Name:           name,
		OwnerId:        ownerId,
		OrganizationId: organizationId,
		Memberships:    memberships,
	}, nil
}

func (lr *LeagueRepository) Save(league *domain.League) error {

	updateID := *league.Id
	if updateID == "" {
		if err := insertIntoLeagues(league); err != nil {
			return err
		}
	} else {
		if err := updateLeague(league); err != nil {
			return err
		}
	}

	return nil
}

func insertIntoLeagues(league *domain.League) error {
	connection := database.GetConnection()

	sql := `INSERT INTO leagues (name, user_id, organization_id) VALUES ($1, $2, $3);`

	_, err := connection.Exec(context.Background(), sql, league.Name, league.OwnerId, league.OrganizationId)

	if err != nil {
		return err
	}

	return nil
}

func updateLeague(league *domain.League) error {
	connection := database.GetConnection()

	sql := `UPDATE leagues SET name=$1, user_id=$2 WHERE id = $3;`

	_, err := connection.Exec(context.Background(), sql, league.Name, league.OwnerId, *league.Id)

	if err != nil {
		return err
	}

	if len(league.Memberships) > 0 {
		var totalValues = len(league.Memberships)
		var values = make([]string, totalValues)

		index := 0
		for _, membership := range league.Memberships {
			values[index] = fmt.Sprintf("('%s', '%s', '%s')", membership.ID, *league.Id, membership.TeamID)
			index++
		}

		sql = `INSERT INTO league_teams (id, league_id, team_id) VALUES ` +
			strings.Join(values, ",") +
			` ON CONFLICT (league_id, team_id) DO NOTHING`

		_, err = connection.Exec(context.Background(), sql)

		if err != nil {
			return err
		}
	}

	membershipValues := []interface{}{*league.Id}
	placeholders := []string{}
	if len(league.Memberships) > 0 {
		for _, membership := range league.Memberships {
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(membershipValues)+1))
			membershipValues = append(membershipValues, membership.ID)
		}
		sql = `DELETE FROM league_teams WHERE league_id=$1 AND id NOT IN (` + strings.Join(placeholders, ",") + `);`
	} else {
		sql = `DELETE FROM league_teams WHERE league_id=$1`
	}

	_, err = connection.Exec(context.Background(), sql, membershipValues...)

	if err != nil {
		return err
	}

	return nil
}

func (lr *LeagueRepository) FetchAll(organizationId string) ([]domain.League, error) {
	connection := database.GetConnection()

	sql := "SELECT " +
		"league_management.leagues.id," +
		"league_management.leagues.name," +
		"league_management.leagues.user_id," +
		"league_management.leagues.organization_id " +
		"FROM league_management.leagues " +
		"JOIN league_management.organizations ON league_management.organizations.id = league_management.leagues.organization_id " +
		"WHERE league_management.organizations.id =$1"

	rows, err := connection.Query(context.Background(), sql, organizationId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var leagues []domain.League

	for rows.Next() {
		var id, name, userId, organizationId string

		if err := rows.Scan(&id, &name, &userId, &organizationId); err != nil {
			return nil, err
		}

		league := domain.NewLeague(id, name, userId, organizationId)

		leagues = append(leagues, league)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return leagues, nil
}

func (lr *LeagueRepository) FetchLeagueMembers(leagueId string) ([]interface{}, error) {
	connection := database.GetConnection()

	sql := `
		SELECT
		    league_teams.id as membership_id,
			teams.id as team_id,
			teams.name as team_name,
			leagues.id as league_id
		FROM league_teams
		JOIN teams ON teams.id = league_teams.team_id
		JOIN leagues ON leagues.id = league_teams.league_id
		WHERE leagues.id=$1;`

	rows, err := connection.Query(context.Background(), sql, leagueId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []interface{}{}

	for rows.Next() {
		var membershipId, teamId, teamName, leagueId string
		if err := rows.Scan(&membershipId, &teamId, &teamName, &leagueId); err != nil {
			return nil, err
		}
		result := make(map[string]string)

		result["membership_id"] = membershipId
		result["team_id"] = teamId
		result["team_name"] = teamName
		result["league_id"] = leagueId

		results = append(results, result)
	}

	return results, nil
}
