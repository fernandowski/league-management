package repositories

import (
	"context"
	"fmt"
	"league-management/internal/organization_management/domain"
	"league-management/internal/shared/database"
	"league-management/internal/shared/dtos"
	"log"
	"strings"
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
		log.Print("error here")
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
			homeId := "bye"
			awayId := "bye"

			if homeTeamID != nil {
				homeId = *homeTeamID
			}
			if awayTeamID != nil {
				awayId = *awayTeamID
			}
			newMatch, _ := domain.NewMatch(homeId, awayId)
			if exists {
				matchesMap[*round] = append(matches, newMatch)
			} else {
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

	sql := `INSERT INTO seasons (id, name, league_id, status) VALUES ($1, $2, $3, $4)
            ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name, league_id=EXCLUDED.league_id, status=EXCLUDED.status;`

	_, err := connection.Exec(context.Background(), sql, season.ID, season.Name, season.LeagueId, season.Status)
	if err != nil {
		return err
	}

	values := []string{}
	params := []interface{}{}
	paramIndex := 1

	if len(season.Rounds) > 0 {
		matchSQL := `INSERT INTO season_schedules 
                     (id, season_id, league_id, round, home_team_id, away_team_id, home_team_score, away_team_score) 
                     VALUES %s 
                     ON CONFLICT (id) DO UPDATE 
                     SET league_id = EXCLUDED.league_id, 
                         round = EXCLUDED.round, 
                         home_team_id = EXCLUDED.home_team_id, 
                         away_team_id = EXCLUDED.away_team_id, 
                         home_team_score = EXCLUDED.home_team_score, 
                         away_team_score = EXCLUDED.away_team_score;`

		for _, round := range season.Rounds {
			for _, match := range round.Matches {

				parametrizedValues := fmt.Sprintf(
					"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
					paramIndex, paramIndex+1, paramIndex+2, paramIndex+3, paramIndex+4, paramIndex+5, paramIndex+6, paramIndex+7,
				)
				values = append(values, parametrizedValues)

				params = append(
					params,
					match.ID,
					season.ID,
					season.LeagueId,
					round.RoundNumber,
					match.GetHomeTeam(),
					match.GetAwayTeam(),
					match.HomeTeamScore,
					match.AwayTeamScore,
				)
				paramIndex += 8
			}
		}

		query := fmt.Sprintf(matchSQL, strings.Join(values, ", "))
		_, err = connection.Exec(context.Background(), query, params...)
		if err != nil {
			return err
		}
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

func (sr *SeasonRepository) FetchDetails(seasonID string) (map[string]interface{}, error) {
	connection := database.GetConnection()
	sql := `SELECT
				seasons.id as season_id,
				seasons.name as season_name,
				season_schedules.id as match_id,
				season_schedules.round as round,
				away_team.name as away_team_name,
			   	home_team.name as home_team_name
			FROM seasons
			LEFT JOIN season_schedules on seasons.id = season_schedules.season_id
			LEFT JOIN teams away_team on  season_schedules.away_team_id = away_team.id
			LEFT JOIN teams home_team on season_schedules.home_team_id = home_team.id
			WHERE seasons.id =$1 AND seasons.status='planned'::season_status;`

	rows, err := connection.Query(context.Background(), sql, seasonID)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var seasonName, seasonId string
	var round *string
	var matchID, awayTeam, homeTeam *string

	matchesMap := make(map[string][]interface{})

	for rows.Next() {

		if err := rows.Scan(&seasonId, &seasonName, &matchID, &round, &awayTeam, &homeTeam); err != nil {
			return nil, err
		}

		if matchID != nil {
			matches, exists := matchesMap[*round]
			homeTeamName := "bye"
			awayTeamName := "bye"

			if homeTeam != nil {
				homeTeamName = *homeTeam
			}
			if awayTeam != nil {
				awayTeamName = *awayTeam
			}

			newMatch := map[string]interface{}{
				"id":        matchID,
				"home_team": homeTeamName,
				"away_team": awayTeamName,
			}
			if exists {
				matchesMap[*round] = append(matches, newMatch)
			} else {
				matchesMap[*round] = []interface{}{newMatch}
			}
		}
	}

	result := map[string]interface{}{
		"id":     seasonId,
		"name":   seasonName,
		"rounds": map[string]interface{}{},
	}
	var rounds []interface{}

	for key, value := range matchesMap {
		round := map[string]interface{}{
			"round_number": key,
			"matches":      value,
		}
		rounds = append(rounds, round)
	}

	result["rounds"] = rounds

	return result, nil
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
