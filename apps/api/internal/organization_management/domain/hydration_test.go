package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRehydrateRound_DefensivelyCopiesMatches(t *testing.T) {
	matches := []Match{
		RehydrateMatch(MatchState{
			ID:         "match-1",
			HomeTeamID: "team-1",
			AwayTeamID: "team-2",
			Status:     MatchStatusInProgress,
		}),
	}

	round := RehydrateRound(1, matches)
	matches[0].homeTeamID = "mutated-home"
	snapshot := round.Snapshot()

	assert.Equal(t, 1, snapshot.RoundNumber)
	assert.Len(t, snapshot.Matches, 1)
	assert.Equal(t, "team-1", snapshot.Matches[0].HomeTeamID)
}

func TestRehydrateMatch_ReconstructsPersistedMatchState(t *testing.T) {
	match := RehydrateMatch(MatchState{
		ID:            "match-1",
		HomeTeamID:    "team-1",
		AwayTeamID:    "team-2",
		HomeTeamScore: 3,
		AwayTeamScore: 2,
		Status:        MatchStatusFinished,
		RefereeID:     "ref-1",
		PlayoffTieID:  "tie-1",
		MatchOrder:    2,
	})

	snapshot := match.Snapshot()
	assert.Equal(t, "match-1", snapshot.ID)
	assert.Equal(t, "team-1", snapshot.HomeTeamID)
	assert.Equal(t, "team-2", snapshot.AwayTeamID)
	assert.Equal(t, 3, snapshot.HomeTeamScore)
	assert.Equal(t, 2, snapshot.AwayTeamScore)
	assert.Equal(t, MatchStatusFinished, snapshot.Status)
	assert.Equal(t, "ref-1", snapshot.RefereeID)
	assert.Equal(t, "tie-1", snapshot.PlayoffTieID)
	assert.Equal(t, 2, snapshot.MatchOrder)
}

func TestRehydratePlayoffBracket_DefensivelyCopiesNestedState(t *testing.T) {
	winnerTeamID := "team-1"
	rounds := []PlayoffBracketRoundSnapshot{
		{
			Name:  "final",
			Order: 1,
			Ties: []PlayoffTieSnapshot{
				{
					ID:           "tie-1",
					RoundName:    "final",
					RoundOrder:   1,
					SlotOrder:    1,
					HomeTeamID:   "team-1",
					AwayTeamID:   "team-2",
					WinnerTeamID: &winnerTeamID,
					Matches: []MatchSnapshot{
						{
							ID:         "match-1",
							HomeTeamID: "team-1",
							AwayTeamID: "team-2",
							Status:     MatchStatusScheduled,
						},
					},
				},
			},
		},
	}

	bracket := RehydratePlayoffBracket(rounds)
	rounds[0].Ties[0].Matches[0].AwayTeamID = "mutated-away"
	winnerTeamID = "mutated-winner"

	assert.NotNil(t, bracket)
	assert.Equal(t, "team-2", bracket.Rounds[0].Ties[0].Matches[0].AwayTeamID)
	assert.NotNil(t, bracket.Rounds[0].Ties[0].WinnerTeamID)
	assert.Equal(t, "team-1", *bracket.Rounds[0].Ties[0].WinnerTeamID)
}

func TestSeasonAccessors_ReturnDefensiveCopies(t *testing.T) {
	championTeamID := "team-1"
	season := RehydrateSeasonFromSnapshot(SeasonSnapshot{
		ID:       "season-1",
		LeagueID: "league-1",
		Name:     "Spring",
		Status:   SeasonStatusInProgress,
		Phase:    SeasonPhasePlayoffs,
		Version:  2,
		Rounds: []RoundSnapshot{
			{
				RoundNumber: 1,
				Matches: []MatchSnapshot{
					{
						ID:         "match-1",
						HomeTeamID: "team-1",
						AwayTeamID: "team-2",
						Status:     MatchStatusInProgress,
					},
				},
			},
		},
		PlayoffRules: &PlayoffRulesSnapshot{
			QualificationType: "top_n",
			QualifierCount:    2,
			Rounds:            []PlayoffRoundRuleSnapshot{{Name: "final", Legs: 1, TiedAggregateResolution: "clear_winner_required"}},
		},
		PlayoffBracket: RehydratePlayoffBracket([]PlayoffBracketRoundSnapshot{
			{
				Name:  "final",
				Order: 1,
				Ties: []PlayoffTieSnapshot{
					{
						ID:           "tie-1",
						HomeTeamID:   "team-1",
						AwayTeamID:   "team-2",
						WinnerTeamID: &championTeamID,
						Matches: []MatchSnapshot{
							{
								ID:         "match-1",
								HomeTeamID: "team-1",
								AwayTeamID: "team-2",
								Status:     MatchStatusFinished,
							},
						},
					},
				},
			},
		}),
		ChampionTeamID: &championTeamID,
	})

	rounds := season.PlannedRounds()
	rules := season.Rules()
	bracket := season.Bracket()
	champion := season.ChampionTeam()
	snapshot := season.Snapshot()

	rounds[0].Matches[0].HomeTeamID = "mutated-home"
	rules.Rounds[0].Name = "mutated-round"
	bracket.Rounds[0].Ties[0].Matches[0].AwayTeamID = "mutated-away"
	*champion = "mutated-champion"
	snapshot.Rounds[0].Matches[0].HomeTeamID = "mutated-snapshot"

	assert.Equal(t, "season-1", season.SeasonID())
	assert.Equal(t, "league-1", season.LeagueID())
	assert.Equal(t, "Spring", season.SeasonName())
	assert.Equal(t, SeasonStatusInProgress, season.CurrentStatus())
	assert.Equal(t, SeasonPhasePlayoffs, season.CurrentPhase())
	assert.Equal(t, 2, season.CurrentVersion())
	matchSnapshot, ok := season.FindMatchSnapshot("match-1")
	assert.True(t, ok)
	assert.Equal(t, "team-1", matchSnapshot.HomeTeamID)
	assert.Equal(t, "final", season.PlayoffRoundRules()[0].Name)
	tie, ok := season.FindPlayoffTie("tie-1")
	assert.True(t, ok)
	assert.Equal(t, "team-2", tie.Matches[0].AwayTeamID)
	assert.Equal(t, "team-1", *season.ChampionTeam())
	assert.Equal(t, "team-1", season.Snapshot().Rounds[0].Matches[0].HomeTeamID)
}
