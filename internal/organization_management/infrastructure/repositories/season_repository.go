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

func (sr *SeasonRepository) FindByID(seasonID string) (*domain.Season, error) {
	connection := database.GetConnection()

	sql := `SELECT
				seasons.id as season_id,
				seasons.name as season_name,
				seasons.league_id as league_id,
				seasons.status as season_status,
				season_schedules.id as match_id,
				season_schedules.round as round,
				season_schedules.home_team_id as home_team_id,
				season_schedules.away_team_id as away_team_id,
				season_schedules.home_team_score as home_team_score,
				season_schedules.away_team_score as away_team_score
			FROM seasons
			LEFT JOIN season_schedules ON season_schedules.season_id = seasons.id
			WHERE seasons.id=$1`

	rows, err := connection.Query(context.Background(), sql, seasonID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	matchesMap := make(map[int][]domain.Match)

	var seasonId, seasonName, leagueId, seasonStatus string
	var matchID, homeTeamID, awayTeamID *string
	var round, homeTeamScore, awayTeamScore *int

	for rows.Next() {
		if err := rows.Scan(&seasonId, &seasonName, &leagueId, &seasonStatus, &matchID, &round, &homeTeamID, &awayTeamID, &homeTeamScore, &awayTeamScore); err != nil {
			return nil, err
		}

		if matchID != nil {
			matches, exists := matchesMap[*round]
			if exists {
				newMatch, _ := domain.NewMatch(*homeTeamID, *awayTeamID)
				matchesMap[*round] = append(matches, newMatch)
			} else {
				newMatch, _ := domain.NewMatch(*homeTeamID, *awayTeamID)
				matchesMap[*round] = []domain.Match{newMatch}
			}
		}
	}

	rounds := []domain.Round{}

	for key, value := range matchesMap {
		newRound := domain.NewRound(key)
		round := newRound.AddMatches(value)
		rounds = append(rounds, *round)
	}

	season := domain.Season{
		ID:             seasonId,
		LeagueId:       leagueId,
		Name:           seasonName,
		Status:         domain.SeasonStatus(seasonStatus),
		Rounds:         rounds,
		MatchLocations: nil,
	}

	return &season, nil
}

func (sr *SeasonRepository) Save(season *domain.Season) error {
	connection := database.GetConnection()

	sql := `INSERT INTO seasons (name, league_id, status) VALUES ($1, $2, $3);`

	_, err := connection.Exec(context.Background(), sql, season.Name, season.LeagueId, season.Status)
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
