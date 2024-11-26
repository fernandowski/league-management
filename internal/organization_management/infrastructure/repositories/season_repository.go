package repositories

import (
	"context"
	"league-management/internal/organization_management/domain"
	"league-management/internal/shared/database"
)

type SeasonRepository struct {
}

func NewSeasonRepository() *SeasonRepository {
	return &SeasonRepository{}
}

func (sr *SeasonRepository) Save(season *domain.Season) error {
	connection := database.GetConnection()

	sql := `INSERT INTO seasons (name, league_id) VALUES ($1, $2);`

	_, err := connection.Exec(context.Background(), sql, season.Name, season.LeagueId)
	if err != nil {
		return err
	}

	return nil
}
