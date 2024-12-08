package repositories

import (
	"context"
	"fmt"
	"league-management/internal/organization_management/domain"
	"league-management/internal/shared/database"
	"league-management/internal/shared/dtos"
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

func (sr *SeasonRepository) Search(orgOwnerID, leagueID string, searchDTO dtos.SearchSeasonDTO) ([]interface{}, int) {
	connection := database.GetConnection()
	queryBuilder := QueryBuilder{}

	sql := `SELECT seasons.id, seasons.name, seasons.league_id
			FROM organizations
			JOIN leagues ON leagues.organization_id=organizations.id
			JOIN seasons ON seasons.league_id=leagues.id`

	queryBuilder.SetQuery(sql)

	queryBuilder.AddFilter("organizations.user_id=$"+fmt.Sprint(len(queryBuilder.parameters)+1), orgOwnerID)
	queryBuilder.AddFilter("leagues.id=$"+fmt.Sprint(len(queryBuilder.parameters)+1), leagueID)

	if searchDTO.Term != "" {
		queryBuilder.AddFilter("seasons.name=$"+fmt.Sprint(len(queryBuilder.parameters)+1), searchDTO.Term)
	}

	queryBuilder.SetPagination(searchDTO.Limit, searchDTO.Offset)

	query, parameters := queryBuilder.BuildQuery()

	rows, err := connection.Query(context.Background(), query, parameters...)

	var results = []interface{}{}
	if err != nil {
		return results, 0
	}

	for rows.Next() {
		var seasonId, seasonName, leagueId string
		if err := rows.Scan(&seasonId, &seasonName, &leagueId); err != nil {
			return results, 0
		}
		season := make(map[string]string)

		season["id"] = seasonId
		season["name"] = seasonName
		season["league_id"] = leagueId

		results = append(results, season)
	}

	total := getCount(orgOwnerID, leagueID, searchDTO)
	return results, total
}

func getCount(orgOwnerID, leagueID string, searchDTO dtos.SearchSeasonDTO) int {
	connection := database.GetConnection()
	queryBuilder := QueryBuilder{}

	sql := `SELECT COUNT(*) AS total
			FROM organizations
			JOIN leagues ON leagues.organization_id=organizations.id
			JOIN seasons ON seasons.league_id=leagues.id`

	queryBuilder.SetQuery(sql)

	queryBuilder.AddFilter("organizations.user_id=$"+fmt.Sprint(len(queryBuilder.parameters)+1), orgOwnerID)
	queryBuilder.AddFilter("leagues.id=$"+fmt.Sprint(len(queryBuilder.parameters)+1), leagueID)

	if searchDTO.Term != "" {
		queryBuilder.AddFilter("seasons.name=$"+fmt.Sprint(len(queryBuilder.parameters)+1), searchDTO.Term)
	}

	query, parameters := queryBuilder.BuildQuery()

	row := connection.QueryRow(context.Background(), query, parameters...)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0
	}

	return total
}
