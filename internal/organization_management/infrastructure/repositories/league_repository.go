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
		Id:             id,
		Name:           name,
		OwnerId:        ownerId,
		OrganizationId: organizationId,
		TeamIds:        []string{},
	}, nil
}

func (lr *LeagueRepository) Save(league *domain.League) error {
	connection := database.GetConnection()

	sql := "INSERT INTO league_management.leagues (name, user_id, organization_id) VALUES ($1, $2, $3);"

	_, err := connection.Exec(context.Background(), sql, league.Name, league.OwnerId, league.OrganizationId)

	if err != nil {
		panic(err.Error())
	}

	return err
}
