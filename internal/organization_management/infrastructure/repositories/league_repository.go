package repositories

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"league-management/internal/organization_management/domain"
	"league-management/internal/shared/database"
	"time"
)

type LeagueRepository struct {
}

func NewLeagueRepository() LeagueRepository {
	return LeagueRepository{}
}

func (lr *LeagueRepository) FindById(leagueId string) (*domain.League, error) {
	connection := database.GetConnection()

	sql := "SELECT id, name, user_id, organization_id, created_at, updated_at" +
		"FROM league_management.leagues" +
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
		panic(err)
	}

	return &domain.League{
		Id:             nil,
		Name:           name,
		OwnerId:        ownerId,
		OrganizationId: organizationId,
		Memberships:    []domain.LeagueMembership{},
	}, nil
}

func (lr *LeagueRepository) Save(league *domain.League) error {

	updateID := league.Id

	if updateID == nil {
		if err := insertIntoLeagues(league); err != nil {
			return err
		}
	} else {

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

	sql := `UPDATE leagues SET name=$1, user_id=$2, organization_id=$3 WHERE id = $4;`

	_, err := connection.Exec(context.Background(), sql, league.Name, league.OwnerId, league.OrganizationId, *league.Id)

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
