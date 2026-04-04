package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"league-management/internal/organization_management/domain"
	"league-management/internal/shared/app_errors"
	"league-management/internal/shared/database"
	"league-management/internal/shared/dtos"
	"strings"
	"time"
)

type LeagueRepository struct {
}

func NewLeagueRepository() *LeagueRepository {
	return &LeagueRepository{}
}

func (lr *LeagueRepository) FindById(leagueId string) (*domain.League, error) {
	return lr.findByIdTx(context.Background(), database.GetConnection(), leagueId, false)
}

func (lr *LeagueRepository) findByIdTx(ctx context.Context, querier interface {
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
}, leagueId string, forUpdate bool) (*domain.League, error) {
	sql := `SELECT 
    		leagues.id,
    		leagues.name,
    		leagues.user_id,
    		leagues.organization_id,
    		leagues.created_at,
    		leagues.updated_at,
    		seasons.id as season_status
		FROM league_management.leagues
		LEFT JOIN league_management.seasons
			ON seasons.league_id=leagues.id
			AND seasons.status IN ('pending', 'planned', 'in_progress', 'paused')
		WHERE leagues.id=$1`
	if forUpdate {
		sql += ` FOR UPDATE`
	}

	var id string
	var name string
	var ownerId string
	var organizationId string
	var seasonId *string
	var dateCreated time.Time
	var dateUpdated time.Time

	err := querier.QueryRow(ctx, sql, leagueId).Scan(&id, &name, &ownerId, &organizationId, &dateCreated, &dateUpdated, &seasonId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	sql = `SELECT id, team_id FROM league_management.league_teams WHERE league_id=$1 ORDER BY id`

	rows, err := querier.Query(ctx, sql, leagueId)
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

	league := domain.League{
		Id:             &id,
		Name:           name,
		OwnerId:        ownerId,
		OrganizationId: organizationId,
		Memberships:    memberships,
		ActiveSeason:   "",
	}

	if seasonId != nil {
		league.ActiveSeason = *seasonId
	}

	return &league, nil
}

func (lr *LeagueRepository) WithLockedLeague(leagueId string, fn func(league *domain.League) error) error {
	return database.WithTx(context.Background(), func(tx pgx.Tx) error {
		league, err := lr.findByIdTx(context.Background(), tx, leagueId, true)
		if err != nil {
			return err
		}

		if err := fn(league); err != nil {
			return err
		}

		return lr.saveTx(context.Background(), tx, league)
	})
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return app_errors.ErrDuplicateResource
		}
		return err
	}

	return nil
}

func updateLeague(league *domain.League) error {
	return (&LeagueRepository{}).saveTx(context.Background(), database.GetConnection(), league)
}

func (lr *LeagueRepository) saveTx(ctx context.Context, tx interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
}, league *domain.League) error {
	sql := `UPDATE league_management.leagues SET name=$1, user_id=$2 WHERE id = $3;`

	_, err := tx.Exec(ctx, sql, league.Name, league.OwnerId, *league.Id)
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

		sql = `INSERT INTO league_management.league_teams (id, league_id, team_id) VALUES ` +
			strings.Join(values, ",") +
			` ON CONFLICT (league_id, team_id) DO NOTHING`

		_, err = tx.Exec(ctx, sql)

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
		sql = `DELETE FROM league_management.league_teams WHERE league_id=$1 AND id NOT IN (` + strings.Join(placeholders, ",") + `);`
	} else {
		sql = `DELETE FROM league_management.league_teams WHERE league_id=$1`
	}

	_, err = tx.Exec(ctx, sql, membershipValues...)

	if err != nil {
		return err
	}

	return nil
}

func (lr *LeagueRepository) Search(searchDTO dtos.LeagueSearchDTO) (*map[string]interface{}, error) {
	connection := database.GetConnection()
	queryBuilder := QueryBuilder{}

	sql := "SELECT " +
		"league_management.leagues.id, " +
		"league_management.leagues.name " +
		"FROM league_management.leagues " +
		"JOIN league_management.organizations ON league_management.organizations.id = league_management.leagues.organization_id "

	queryBuilder.SetQuery(sql)

	queryBuilder.AddFilter("league_management.organizations.id=$"+fmt.Sprint(len(queryBuilder.parameters)+1), searchDTO.OrganizationID)
	queryBuilder.SetPagination(searchDTO.Limit, searchDTO.Offset)

	sql, parameters := queryBuilder.BuildQuery()
	rows, err := connection.Query(context.Background(), sql, parameters...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var leagues []interface{}

	leagueMemberCount, err := organizationLeagueMembersCount(searchDTO.OrganizationID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id, name string

		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}

		league := make(map[string]interface{})
		league["id"] = id
		league["name"] = name
		league["team_ids"] = []interface{}{}
		league["total_members"] = leagueMemberCount[id]

		leagues = append(leagues, league)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	total, err := searchCount(searchDTO)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	result["data"] = leagues
	result["total"] = total
	return &result, nil
}

func searchCount(searchDTO dtos.LeagueSearchDTO) (int, error) {
	connection := database.GetConnection()
	queryBuilder := QueryBuilder{}

	sql := "SELECT " +
		"COUNT(*) AS total " +
		"FROM league_management.leagues " +
		"JOIN league_management.organizations ON league_management.organizations.id = league_management.leagues.organization_id "

	queryBuilder.SetQuery(sql)

	queryBuilder.AddFilter("league_management.organizations.id=$"+fmt.Sprint(len(queryBuilder.parameters)+1), searchDTO.OrganizationID)

	sql, parameters := queryBuilder.BuildQuery()
	row := connection.QueryRow(context.Background(), sql, parameters...)

	var count int

	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func organizationLeagueMembersCount(organizationId string) (map[string]int, error) {
	connection := database.GetConnection()

	sql := `SELECT leagues.id, organizations.id, count(*) as total_members
			FROM organizations
			JOIN leagues ON leagues.organization_id = organizations.id
			JOIN league_teams ON league_teams.league_id = leagues.id
			WHERE organizations.id=$1
			GROUP BY 1, 2`

	leagueTeamCount := make(map[string]int)

	rows, err := connection.Query(context.Background(), sql, organizationId)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var leagueId, organizationId string
		var totalMembers int

		if err := rows.Scan(&leagueId, &organizationId, &totalMembers); err != nil {
			return nil, err
		}

		leagueTeamCount[leagueId] = totalMembers
	}

	return leagueTeamCount, nil
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

func (lr *LeagueRepository) FetchLeagueDetails(league domain.League) (map[string]interface{}, error) {
	connection := database.GetConnection()

	sql := `SELECT 
    			leagues.id, 
    			leagues.name, 
    			seasons.name as season_name, 
    			seasons.id as season_id, 
    			seasons.status as season_status
			FROM leagues 
			LEFT JOIN seasons ON seasons.league_id=leagues.id AND seasons.status in ('pending', 'planned', 'in_progress', 'paused')
			WHERE leagues.id=$1;`

	var leagueId, leagueName string
	var seasonId, seasonName, seasonStatus *string

	err := connection.QueryRow(context.Background(), sql, *league.Id).Scan(&leagueId, &leagueName, &seasonName, &seasonId, &seasonStatus)

	if err != nil {
		return nil, err
	}

	leagueMembershipCount, err := organizationLeagueMembersCount(league.OrganizationId)
	if err != nil {
		return nil, err
	}

	leagueProjection := make(map[string]interface{})
	leagueProjection["name"] = leagueName
	leagueProjection["id"] = leagueId
	leagueProjection["active_members"] = leagueMembershipCount[*league.Id]
	leagueProjection["season"] = nil

	if seasonId != nil {
		seasonProjection := make(map[string]string)
		seasonProjection["id"] = *seasonId
		seasonProjection["name"] = *seasonName
		seasonProjection["status"] = *seasonStatus
		leagueProjection["season"] = seasonProjection
	}

	return leagueProjection, nil
}
