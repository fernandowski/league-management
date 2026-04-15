package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"league-management/internal/organization_management/domain"
	"league-management/internal/shared/app_errors"
	"league-management/internal/shared/database"
	"league-management/internal/shared/dtos"
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
				seasons.season_phase as season_phase,
				seasons.version as version,
				seasons.champion_team_id as champion_team_id,
				season_schedules.id as match_id,
				season_schedules.round as round,
				season_schedules.home_team_id as home_team_id,
				season_schedules.away_team_id as away_team_id,
				season_schedules.home_team_score as home_team_score,
				season_schedules.away_team_score as away_team_score,
				season_schedules.status as match_status,
				season_schedules.referee_id as referee_id
			FROM seasons
			LEFT JOIN season_schedules ON season_schedules.season_id = seasons.id
			WHERE seasons.id=$1`

	rows, err := connection.Query(context.Background(), sql, seasonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matchesMap := make(map[int][]domain.Match)

	var seasonId, seasonName, leagueId, seasonStatus, seasonPhase string
	var version int
	var championTeamID *string
	var matchID, homeTeamID, awayTeamID, refereeID, matchStatus *string
	var round, homeTeamScore, awayTeamScore *int

	foundRows := false

	for rows.Next() {
		if err := rows.Scan(&seasonId, &seasonName, &leagueId, &seasonStatus, &seasonPhase, &version, &championTeamID, &matchID, &round, &homeTeamID, &awayTeamID, &homeTeamScore, &awayTeamScore, &matchStatus, &refereeID); err != nil {
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
			newMatch, _ := domain.NewMatch(matchID, homeId, awayId)
			newMatch.AwayTeamScore = *awayTeamScore
			newMatch.HomeTeamScore = *homeTeamScore
			if matchStatus != nil {
				newMatch.Status = domain.MatchStatus(*matchStatus)
			}
			if refereeID != nil {
				newMatch.RefereeID = *refereeID
			}
			if exists {
				matchesMap[*round] = append(matches, newMatch)
			} else {
				matchesMap[*round] = []domain.Match{newMatch}
			}
		}
		foundRows = true
	}

	if !foundRows {
		return nil, errors.New("no season found")
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
		Phase:          domain.SeasonPhase(seasonPhase),
		Version:        version,
		Rounds:         rounds,
		MatchLocations: nil,
		ChampionTeamID: championTeamID,
	}

	playoffRules, err := sr.findPlayoffRules(seasonID)
	if err != nil {
		return nil, err
	}
	season.PlayoffRules = playoffRules
	playoffBracket, err := sr.findPlayoffBracket(seasonID)
	if err != nil {
		return nil, err
	}
	season.PlayoffBracket = playoffBracket

	return &season, nil
}

func (sr *SeasonRepository) Save(season *domain.Season) error {
	return database.WithTx(context.Background(), func(tx pgx.Tx) error {
		if season.Version <= 0 {
			sql := `INSERT INTO league_management.seasons (id, name, league_id, status, season_phase, version, champion_team_id) VALUES ($1, $2, $3, $4, $5, $6, $7)`
			_, err := tx.Exec(context.Background(), sql, season.ID, season.Name, season.LeagueId, season.Status, season.Phase, 1, season.ChampionTeamID)
			if err != nil {
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) && pgErr.Code == "23505" {
					return app_errors.ErrDuplicateResource
				}
				return err
			}
			season.Version = 1
		} else {
			sql := `UPDATE league_management.seasons
				SET name=$1, league_id=$2, status=$3, season_phase=$4, champion_team_id=$5, version=version+1
				WHERE id=$6 AND version=$7`

			tag, err := tx.Exec(context.Background(), sql, season.Name, season.LeagueId, season.Status, season.Phase, season.ChampionTeamID, season.ID, season.Version)
			if err != nil {
				return err
			}
			if tag.RowsAffected() == 0 {
				return app_errors.ErrConcurrentModification
			}
			season.Version++
		}

		if err := sr.savePlayoffRules(tx, season); err != nil {
			return err
		}
		if err := sr.savePlayoffBracket(tx, season); err != nil {
			return err
		}

		values := []string{}
		params := []interface{}{}
		paramIndex := 1

		if len(season.Rounds) > 0 {
			matchSQL := `INSERT INTO league_management.season_schedules 
					 (id, season_id, league_id, round, home_team_id, away_team_id, home_team_score, away_team_score, status, referee_id) 
					 VALUES %s 
					 ON CONFLICT (id) DO UPDATE 
					 SET league_id = EXCLUDED.league_id, 
						 round = EXCLUDED.round, 
						 home_team_id = EXCLUDED.home_team_id, 
						 away_team_id = EXCLUDED.away_team_id, 
						 home_team_score = EXCLUDED.home_team_score, 
						 away_team_score = EXCLUDED.away_team_score, 
						 status = EXCLUDED.status,
						 referee_id = EXCLUDED.referee_id;`
			for _, round := range season.Rounds {
				for _, match := range round.Matches {
					parametrizedValues := fmt.Sprintf(
						"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
						paramIndex, paramIndex+1, paramIndex+2, paramIndex+3, paramIndex+4, paramIndex+5, paramIndex+6, paramIndex+7, paramIndex+8, paramIndex+9,
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
						match.Status,
						match.RefereeID,
					)
					paramIndex += 10
				}
			}

			query := fmt.Sprintf(matchSQL, strings.Join(values, ", "))
			_, err := tx.Exec(context.Background(), query, params...)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (sr *SeasonRepository) findPlayoffRules(seasonID string) (*domain.PlayoffRules, error) {
	connection := database.GetConnection()

	rulesSQL := `SELECT qualification_type, qualifier_count, reseed_each_round, third_place_match, allow_admin_seed_override
		FROM league_management.season_playoff_rules
		WHERE season_id=$1`

	var qualificationType string
	var qualifierCount int
	var reseedEachRound, thirdPlaceMatch, allowAdminSeedOverride bool

	err := connection.QueryRow(context.Background(), rulesSQL, seasonID).Scan(
		&qualificationType,
		&qualifierCount,
		&reseedEachRound,
		&thirdPlaceMatch,
		&allowAdminSeedOverride,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	roundsSQL := `SELECT round_name, legs, higher_seed_hosts_second_leg, tied_aggregate_resolution
		FROM league_management.season_playoff_round_rules
		WHERE season_id=$1
		ORDER BY round_order ASC`

	rows, err := connection.Query(context.Background(), roundsSQL, seasonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rounds := []domain.PlayoffRoundRule{}
	for rows.Next() {
		var round domain.PlayoffRoundRule
		if err := rows.Scan(&round.Name, &round.Legs, &round.HigherSeedHostsSecondLeg, &round.TiedAggregateResolution); err != nil {
			return nil, err
		}
		rounds = append(rounds, round)
	}

	return &domain.PlayoffRules{
		QualificationType:      qualificationType,
		QualifierCount:         qualifierCount,
		ReseedEachRound:        reseedEachRound,
		ThirdPlaceMatch:        thirdPlaceMatch,
		AllowAdminSeedOverride: allowAdminSeedOverride,
		Rounds:                 rounds,
	}, nil
}

func (sr *SeasonRepository) savePlayoffRules(tx pgx.Tx, season *domain.Season) error {
	if season.PlayoffRules == nil {
		return nil
	}

	rulesSQL := `INSERT INTO league_management.season_playoff_rules
		(season_id, qualification_type, qualifier_count, reseed_each_round, third_place_match, allow_admin_seed_override)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (season_id) DO UPDATE
		SET qualification_type = EXCLUDED.qualification_type,
			qualifier_count = EXCLUDED.qualifier_count,
			reseed_each_round = EXCLUDED.reseed_each_round,
			third_place_match = EXCLUDED.third_place_match,
			allow_admin_seed_override = EXCLUDED.allow_admin_seed_override`

	_, err := tx.Exec(
		context.Background(),
		rulesSQL,
		season.ID,
		season.PlayoffRules.QualificationType,
		season.PlayoffRules.QualifierCount,
		season.PlayoffRules.ReseedEachRound,
		season.PlayoffRules.ThirdPlaceMatch,
		season.PlayoffRules.AllowAdminSeedOverride,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(), `DELETE FROM league_management.season_playoff_round_rules WHERE season_id=$1`, season.ID)
	if err != nil {
		return err
	}

	for index, round := range season.PlayoffRules.Rounds {
		_, err = tx.Exec(
			context.Background(),
			`INSERT INTO league_management.season_playoff_round_rules
				(season_id, round_order, round_name, legs, higher_seed_hosts_second_leg, tied_aggregate_resolution)
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			season.ID,
			index+1,
			round.Name,
			round.Legs,
			round.HigherSeedHostsSecondLeg,
			round.TiedAggregateResolution,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sr *SeasonRepository) findPlayoffBracket(seasonID string) (*domain.PlayoffBracket, error) {
	connection := database.GetConnection()

	var bracketID string
	err := connection.QueryRow(context.Background(), `SELECT id FROM league_management.season_playoff_brackets WHERE season_id=$1`, seasonID).Scan(&bracketID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	rows, err := connection.Query(context.Background(), `SELECT id, round_name, round_order, slot_order, home_seed, away_seed, home_team_id, away_team_id, status, winner_team_id
		FROM league_management.season_playoff_ties
		WHERE playoff_bracket_id=$1
		ORDER BY round_order ASC, slot_order ASC`, bracketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bracket := &domain.PlayoffBracket{Rounds: []domain.PlayoffBracketRound{}}
	roundMap := map[int]int{}

	for rows.Next() {
		var tie domain.PlayoffTie
		var homeSeed, awaySeed *int
		var homeTeamID, awayTeamID, winnerTeamID *string
		if err := rows.Scan(&tie.ID, &tie.RoundName, &tie.RoundOrder, &tie.SlotOrder, &homeSeed, &awaySeed, &homeTeamID, &awayTeamID, &tie.Status, &winnerTeamID); err != nil {
			return nil, err
		}
		if homeSeed != nil {
			tie.HomeSeed = *homeSeed
		}
		if awaySeed != nil {
			tie.AwaySeed = *awaySeed
		}
		if homeTeamID != nil {
			tie.HomeTeamID = *homeTeamID
		}
		if awayTeamID != nil {
			tie.AwayTeamID = *awayTeamID
		}
		if winnerTeamID != nil {
			tie.WinnerTeamID = winnerTeamID
		}
		matches, err := sr.findPlayoffMatches(tie.ID)
		if err != nil {
			return nil, err
		}
		tie.Matches = matches

		index, exists := roundMap[tie.RoundOrder]
		if !exists {
			bracket.Rounds = append(bracket.Rounds, domain.PlayoffBracketRound{
				Name:  tie.RoundName,
				Order: tie.RoundOrder,
				Ties:  []domain.PlayoffTie{},
			})
			index = len(bracket.Rounds) - 1
			roundMap[tie.RoundOrder] = index
		}
		bracket.Rounds[index].Ties = append(bracket.Rounds[index].Ties, tie)
	}

	return bracket, nil
}

func (sr *SeasonRepository) savePlayoffBracket(tx pgx.Tx, season *domain.Season) error {
	if season.PlayoffBracket == nil {
		return nil
	}

	var bracketID string
	err := tx.QueryRow(context.Background(), `SELECT id FROM league_management.season_playoff_brackets WHERE season_id=$1`, season.ID).Scan(&bracketID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			bracketID = uuid.New().String()
			_, err = tx.Exec(context.Background(), `INSERT INTO league_management.season_playoff_brackets (id, season_id, status) VALUES ($1, $2, $3)`, bracketID, season.ID, "generated")
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	_, err = tx.Exec(context.Background(), `DELETE FROM league_management.season_playoff_ties WHERE playoff_bracket_id=$1`, bracketID)
	if err != nil {
		return err
	}

	for _, round := range season.PlayoffBracket.Rounds {
		for _, tie := range round.Ties {
			var homeSeed interface{}
			var awaySeed interface{}
			var homeTeamID interface{}
			var awayTeamID interface{}
			if tie.HomeSeed > 0 {
				homeSeed = tie.HomeSeed
			}
			if tie.AwaySeed > 0 {
				awaySeed = tie.AwaySeed
			}
			if tie.HomeTeamID != "" {
				homeTeamID = tie.HomeTeamID
			}
			if tie.AwayTeamID != "" {
				awayTeamID = tie.AwayTeamID
			}

			var winnerTeamID interface{}
			if tie.WinnerTeamID != nil {
				winnerTeamID = *tie.WinnerTeamID
			}

			_, err = tx.Exec(context.Background(), `INSERT INTO league_management.season_playoff_ties
				(id, playoff_bracket_id, round_name, round_order, slot_order, home_seed, away_seed, home_team_id, away_team_id, status, winner_team_id)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
				tie.ID, bracketID, tie.RoundName, tie.RoundOrder, tie.SlotOrder, homeSeed, awaySeed, homeTeamID, awayTeamID, tie.Status, winnerTeamID,
			)
			if err != nil {
				return err
			}

			for _, match := range tie.Matches {
				var matchHomeTeamID interface{}
				var matchAwayTeamID interface{}
				if match.HomeTeamID != "" {
					matchHomeTeamID = match.HomeTeamID
				}
				if match.AwayTeamID != "" {
					matchAwayTeamID = match.AwayTeamID
				}

				_, err = tx.Exec(context.Background(), `INSERT INTO league_management.season_playoff_matches
					(id, playoff_tie_id, match_order, home_team_id, away_team_id, home_score, away_score, status)
					VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
					match.ID, tie.ID, match.MatchOrder, matchHomeTeamID, matchAwayTeamID, match.HomeTeamScore, match.AwayTeamScore, match.Status,
				)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (sr *SeasonRepository) findPlayoffMatches(tieID string) ([]domain.Match, error) {
	connection := database.GetConnection()

	rows, err := connection.Query(context.Background(), `SELECT id, match_order, home_team_id, away_team_id, home_score, away_score, status
		FROM league_management.season_playoff_matches
		WHERE playoff_tie_id=$1
		ORDER BY match_order ASC`, tieID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.Match{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	matches := []domain.Match{}
	for rows.Next() {
		var match domain.Match
		var homeTeamID, awayTeamID *string
		if err := rows.Scan(&match.ID, &match.MatchOrder, &homeTeamID, &awayTeamID, &match.HomeTeamScore, &match.AwayTeamScore, &match.Status); err != nil {
			return nil, err
		}
		match.PlayoffTieID = tieID
		if homeTeamID != nil {
			match.HomeTeamID = *homeTeamID
		}
		if awayTeamID != nil {
			match.AwayTeamID = *awayTeamID
		}
		matches = append(matches, match)
	}

	return matches, nil
}

func (sr *SeasonRepository) Search(orgOwnerID, leagueID string, searchDTO dtos.SearchSeasonDTO) ([]interface{}, int) {
	connection := database.GetConnection()
	queryBuilder := QueryBuilder{}

	sql := `SELECT seasons.id, seasons.name, seasons.league_id, seasons.status
			FROM organizations
			JOIN leagues ON leagues.organization_id=organizations.id
			JOIN seasons ON seasons.league_id=leagues.id`

	queryBuilder.SetQuery(sql)

	queryBuilder.AddFilter("organizations.user_id=$"+fmt.Sprint(len(queryBuilder.parameters)+1), orgOwnerID)
	queryBuilder.AddFilter("leagues.id=$"+fmt.Sprint(len(queryBuilder.parameters)+1), leagueID)
	queryBuilder.AddFilter("seasons.status!=$"+fmt.Sprint(len(queryBuilder.parameters)+1), "undefined")

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
		var seasonId, seasonName, leagueId, status string
		if err := rows.Scan(&seasonId, &seasonName, &leagueId, &status); err != nil {
			return results, 0
		}
		season := make(map[string]string)

		season["id"] = seasonId
		season["name"] = seasonName
		season["league_id"] = leagueId
		season["status"] = status

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
				seasons.status as season_status,
				seasons.season_phase as season_phase,
				seasons.champion_team_id as champion_team_id,
				season_schedules.id as match_id,
				season_schedules.round as round,
				away_team.name as away_team_name,
			   	home_team.name as home_team_name
			FROM seasons
			LEFT JOIN season_schedules on seasons.id = season_schedules.season_id
			LEFT JOIN teams away_team on  season_schedules.away_team_id = away_team.id
			LEFT JOIN teams home_team on season_schedules.home_team_id = home_team.id
			WHERE seasons.id =$1`

	rows, err := connection.Query(context.Background(), sql, seasonID)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var seasonName, seasonId, seasonStatus, seasonPhase string
	var championTeamID *string
	var round *string
	var matchID, awayTeam, homeTeam *string

	matchesMap := make(map[string][]interface{})

	foundRows := false
	for rows.Next() {

		if err := rows.Scan(&seasonId, &seasonName, &seasonStatus, &seasonPhase, &championTeamID, &matchID, &round, &awayTeam, &homeTeam); err != nil {
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
		foundRows = true
	}

	if !foundRows {
		return nil, errors.New("season not found")
	}

	result := map[string]interface{}{
		"id":               seasonId,
		"name":             seasonName,
		"status":           seasonStatus,
		"phase":            seasonPhase,
		"champion_team_id": championTeamID,
		"rounds":           map[string]interface{}{},
	}
	rounds := []interface{}{}

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

func (sr *SeasonRepository) FetchSeasonStandings(seasonID string) (map[string]interface{}, error) {
	connection := database.GetConnection()

	sql := `WITH filtered_matches AS (
		SELECT
			home_team_id AS team_id,
			home_team_score AS total_goals,
			(home_team_score > away_team_score AND status='finished') AS win,
			(home_team_score < away_team_score AND status='finished') AS losses,
			(home_team_score = away_team_score AND status='finished') AS tie
		FROM league_management.season_schedules
		WHERE season_id = $1 AND away_team_id IS NOT NULL AND home_team_id IS NOT NULL

    	UNION ALL

    	SELECT
        	away_team_id AS team_id,
        	away_team_score AS total_goals,
        	(away_team_score > home_team_score AND status='finished') AS win,
        	(away_team_score < home_team_score AND status='finished') AS losses,
        	(away_team_score = home_team_score AND status='finished') AS tie
    	FROM league_management.season_schedules
    	WHERE season_id = $1 AND away_team_id IS NOT NULL AND home_team_id IS NOT NULL
	)
	SELECT
		team_id,
		teams.name,
		SUM(case WHEN win OR tie OR losses THEN 1 ELSE 0 END) AS games_played,
		SUM(CASE WHEN win THEN 3 WHEN tie THEN 1 ELSE 0 END) AS total_points,
		SUM(total_goals) AS total_goals,
		SUM(CASE WHEN win THEN 1 ELSE 0 END) AS total_wins,
		SUM(CASE WHEN losses THEN 1 ELSE 0 END) AS total_losses,
		SUM(CASE WHEN tie THEN 1 ELSE 0 END) AS total_ties
	FROM filtered_matches
	JOIN teams ON teams.id=filtered_matches.team_id
	WHERE team_id IS NOT NULL
	GROUP BY 1, 2
	ORDER BY total_points DESC;`

	rows, err := connection.Query(context.Background(), sql, seasonID)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	foundRows := false

	var teamId, teamName string
	var gamesPlayed, totalPoints, totalGoals, totalWins, totalLosses, totalTies int

	standingRecords := []interface{}{}

	for rows.Next() {
		if err := rows.Scan(&teamId, &teamName, &gamesPlayed, &totalPoints, &totalGoals, &totalWins, &totalLosses, &totalTies); err != nil {
			return nil, err
		}

		foundRows = true

		standingRecord := make(map[string]interface{})

		standingRecord["team_id"] = teamId
		standingRecord["team_name"] = teamName
		standingRecord["games_played"] = gamesPlayed
		standingRecord["total_points"] = totalPoints
		standingRecord["total_goals"] = totalGoals
		standingRecord["total_wins"] = totalWins
		standingRecord["total_losses"] = totalLosses
		standingRecord["total_ties"] = totalTies

		standingRecords = append(standingRecords, standingRecord)
	}

	if !foundRows {
		return nil, errors.New("no standings for season")
	}

	result := map[string]interface{}{
		"season_id": seasonID,
		"standings": standingRecords,
	}

	return result, nil
}

func (sr *SeasonRepository) FetchSeasonMatchUps(seasonID string) ([]interface{}, error) {
	connection := database.GetConnection()

	sql := `
		SELECT
			season_schedules.id AS match_id,
			(CASE WHEN home_team.name IS NULL THEN 'Bye' ELSE home_team.name END)  AS home_team_name,
			(CASE WHEN away_team.name IS NULL THEN 'Bye' ELSE away_team.name END) AS away_team_name,
			season_schedules.round AS round,
			season_schedules.home_team_score AS home_team_score,
			season_schedules.away_team_score AS away_team_score,
			season_schedules.status AS match_status
		FROM season_schedules
		LEFT JOIN teams AS home_team ON home_team.id=season_schedules.home_team_id
		LEFT JOIN teams AS away_team ON away_team.id=season_schedules.away_team_id
		WHERE season_id=$1;`

	rows, err := connection.Query(context.Background(), sql, seasonID)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	matchesPerRound := map[int]interface{}{}

	for rows.Next() {

		var matchID, homeTeamName, awayTeamName, matchStatus string
		var round, homeScore, awayScore int
		if err := rows.Scan(&matchID, &homeTeamName, &awayTeamName, &round, &homeScore, &awayScore, &matchStatus); err != nil {
			return nil, err
		}

		newMatch := map[string]interface{}{
			"id":         matchID,
			"home_team":  homeTeamName,
			"away_team":  awayTeamName,
			"home_score": homeScore,
			"away_score": awayScore,
			"round":      round,
			"status":     matchStatus,
		}

		matches, exists := matchesPerRound[round]
		if exists {
			matchesPerRound[round] = append(matches.([]interface{}), newMatch)
		} else {
			matchesPerRound[round] = []interface{}{newMatch}
		}
	}

	result := []interface{}{}
	for key, value := range matchesPerRound {
		result = append(result, map[string]interface{}{
			"round":   key,
			"matches": value,
		})
	}

	return result, nil
}

func (sr *SeasonRepository) FetchPlayoffBracket(seasonID string) (map[string]interface{}, error) {
	connection := database.GetConnection()

	var bracketID string
	err := connection.QueryRow(context.Background(), `SELECT id FROM league_management.season_playoff_brackets WHERE season_id=$1`, seasonID).Scan(&bracketID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return map[string]interface{}{
				"season_id": seasonID,
				"generated": false,
				"rounds":    []interface{}{},
			}, nil
		}
		return nil, err
	}

	rows, err := connection.Query(context.Background(), `
		SELECT
			ties.id,
			ties.round_name,
			ties.round_order,
			ties.slot_order,
			ties.home_seed,
			ties.away_seed,
			ties.status,
			ties.winner_team_id,
			home_team.id,
			COALESCE(home_team.name, 'TBD'),
			away_team.id,
			COALESCE(away_team.name, 'TBD')
		FROM league_management.season_playoff_ties AS ties
		LEFT JOIN league_management.teams AS home_team ON home_team.id = ties.home_team_id
		LEFT JOIN league_management.teams AS away_team ON away_team.id = ties.away_team_id
		WHERE ties.playoff_bracket_id=$1
		ORDER BY ties.round_order ASC, ties.slot_order ASC`, bracketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type roundEntry struct {
		name  string
		order int
		ties  []interface{}
	}

	roundMap := map[int]*roundEntry{}
	roundOrder := []int{}

	for rows.Next() {
		var tieID, roundName, status string
		var roundNumber, slotOrder int
		var homeSeed, awaySeed *int
		var winnerTeamID, homeTeamID, homeTeamName, awayTeamID, awayTeamName *string
		if err := rows.Scan(&tieID, &roundName, &roundNumber, &slotOrder, &homeSeed, &awaySeed, &status, &winnerTeamID, &homeTeamID, &homeTeamName, &awayTeamID, &awayTeamName); err != nil {
			return nil, err
		}

		entry, exists := roundMap[roundNumber]
		if !exists {
			entry = &roundEntry{name: roundName, order: roundNumber, ties: []interface{}{}}
			roundMap[roundNumber] = entry
			roundOrder = append(roundOrder, roundNumber)
		}

		tie := map[string]interface{}{
			"id":         tieID,
			"slot_order": slotOrder,
			"status":     status,
			"home_seed":  homeSeed,
			"away_seed":  awaySeed,
			"home_team": map[string]interface{}{
				"id":   homeTeamID,
				"name": valueOrFallback(homeTeamName, "TBD"),
			},
			"away_team": map[string]interface{}{
				"id":   awayTeamID,
				"name": valueOrFallback(awayTeamName, "TBD"),
			},
			"winner_team_id": winnerTeamID,
			"matches":        []interface{}{},
		}

		matchRows, err := connection.Query(context.Background(), `
			SELECT
				legs.id,
				legs.match_order,
				legs.home_score,
				legs.away_score,
				legs.status,
				COALESCE(home_team.name, 'TBD'),
				COALESCE(away_team.name, 'TBD')
			FROM league_management.season_playoff_matches AS legs
			LEFT JOIN league_management.teams AS home_team ON home_team.id = legs.home_team_id
			LEFT JOIN league_management.teams AS away_team ON away_team.id = legs.away_team_id
			WHERE legs.playoff_tie_id=$1
			ORDER BY legs.match_order ASC`, tieID)
		if err != nil {
			return nil, err
		}

		matches := []interface{}{}
		for matchRows.Next() {
			var matchID, matchStatus, homeTeamName, awayTeamName string
			var matchOrder, homeScore, awayScore int
			if err := matchRows.Scan(&matchID, &matchOrder, &homeScore, &awayScore, &matchStatus, &homeTeamName, &awayTeamName); err != nil {
				matchRows.Close()
				return nil, err
			}
			matches = append(matches, map[string]interface{}{
				"id":          matchID,
				"match_order": matchOrder,
				"home_score":  homeScore,
				"away_score":  awayScore,
				"status":      matchStatus,
				"home_team":   homeTeamName,
				"away_team":   awayTeamName,
			})
		}
		matchRows.Close()
		tie["matches"] = matches

		entry.ties = append(entry.ties, tie)
	}

	resultRounds := []interface{}{}
	for _, number := range roundOrder {
		entry := roundMap[number]
		resultRounds = append(resultRounds, map[string]interface{}{
			"name":  entry.name,
			"order": entry.order,
			"ties":  entry.ties,
		})
	}

	return map[string]interface{}{
		"season_id": seasonID,
		"generated": true,
		"rounds":    resultRounds,
	}, nil
}

func valueOrFallback(value *string, fallback string) string {
	if value == nil || *value == "" {
		return fallback
	}
	return *value
}
