package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSeason_ScheduleRounds(t *testing.T) {
	t.Run("Zero teams should error out", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusPending,
			Rounds:   nil,
		}

		leagueId := "league-id"
		league := League{
			Id:             &leagueId,
			Name:           "test-league",
			OwnerId:        "test-owner-id",
			OrganizationId: "test-org-id",
			Memberships:    []LeagueMembership{},
		}

		err := season.ScheduleRounds(league)
		assert.Error(t, err)
	})

	t.Run("Season with even teams without rules", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusPending,
		}

		member1 := LeagueMembership{
			ID:     "test-id",
			TeamID: "team A",
		}
		member2 := LeagueMembership{
			ID:     "test-id",
			TeamID: "team B",
		}
		member3 := LeagueMembership{
			ID:     "test-id",
			TeamID: "team C",
		}
		member4 := LeagueMembership{
			ID:     "test-id",
			TeamID: "team D",
		}

		leagueId := "league-id"
		league := League{
			Id:             &leagueId,
			Name:           "test-league",
			OwnerId:        "test-owner-id",
			OrganizationId: "test-org-id",
			Memberships:    []LeagueMembership{member1, member2, member3, member4},
		}

		err := season.ScheduleRounds(league)

		firstRound := season.Rounds[0]

		assert.NoError(t, err)
		assert.Len(t, season.Rounds, 3)
		assert.Equal(t, season.Status, SeasonStatusPlanned)
		assert.Equal(t, 1, firstRound.RoundNumber)
		assert.Len(t, firstRound.Matches, 2)
	})

	t.Run("Season with odd teams without rules", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusPending,
		}

		member1 := LeagueMembership{
			ID:     "test-id",
			TeamID: "team A",
		}
		member2 := LeagueMembership{
			ID:     "test-id",
			TeamID: "team B",
		}
		member3 := LeagueMembership{
			ID:     "test-id",
			TeamID: "team C",
		}

		leagueId := "league-id"
		league := League{
			Id:             &leagueId,
			Name:           "test-league",
			OwnerId:        "test-owner-id",
			OrganizationId: "test-org-id",
			Memberships:    []LeagueMembership{member1, member2, member3},
		}

		err := season.ScheduleRounds(league)

		firstRound := season.Rounds[0]

		assert.NoError(t, err)
		assert.Len(t, season.Rounds, 3)
		assert.Equal(t, season.Status, SeasonStatusPlanned)
		assert.Equal(t, 1, firstRound.RoundNumber)
		assert.Len(t, firstRound.Matches, 2)

		foundByeRival := false
		for _, match := range firstRound.Matches {
			if match.AwayTeamID == "bye" || match.HomeTeamID == "bye" {
				foundByeRival = true
			}
		}

		assert.True(t, foundByeRival, "did not find by rival.")
	})

	t.Run("Start season in pending status should fail", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusPending,
		}

		_, err := season.Start()
		assert.Error(t, err)
	})

	t.Run("Start season in_progress status should fail", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusInProgress,
		}

		_, err := season.Start()
		assert.Error(t, err)
	})

	t.Run("Start season paused status should fail", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusPaused,
		}

		_, err := season.Start()
		assert.Error(t, err)
	})

	t.Run("Start season finished status should fail", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusFinished,
		}

		_, err := season.Start()
		assert.Error(t, err)
	})

	t.Run("Start season undefined status should fail", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusUndefined,
		}

		_, err := season.Start()
		assert.Error(t, err)
	})

	t.Run("Start Season with empty rounds should fail", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusPending,
			Rounds:   []Round{},
		}

		_, err := season.Start()
		assert.Error(t, err)
	})

	t.Run("Season with even teams without rules", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusPending,
		}

		member1 := LeagueMembership{
			ID:     "test-id",
			TeamID: "team A",
		}
		member2 := LeagueMembership{
			ID:     "test-id",
			TeamID: "team B",
		}
		member3 := LeagueMembership{
			ID:     "test-id",
			TeamID: "team C",
		}
		member4 := LeagueMembership{
			ID:     "test-id",
			TeamID: "team D",
		}

		leagueId := "league-id"
		league := League{
			Id:             &leagueId,
			Name:           "test-league",
			OwnerId:        "test-owner-id",
			OrganizationId: "test-org-id",
			Memberships:    []LeagueMembership{member1, member2, member3, member4},
		}

		_ = season.ScheduleRounds(league)
		newSeason, _ := season.Start()

		assert.Equal(t, SeasonStatusInProgress, newSeason.Status)
		assert.Equal(t, SeasonStatusPlanned, season.Status)
		for _, match := range newSeason.Rounds[0].Matches {
			assert.Equal(t, MatchStatusInProgress, match.Status)
		}
	})

	t.Run("Change game score to negative number should fail", func(t *testing.T) {

		testId1 := "test_id_2"
		match1, _ := NewMatch(&testId1, "one", "two")

		testId2 := "test_id_2"
		match2, _ := NewMatch(&testId2, "one", "two")

		round := NewRound(1)
		round.Matches = []Match{match1, match2}

		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Rounds:   []Round{round},
			Status:   SeasonStatusInProgress,
		}

		_, err := season.ChangeMatchScore(testId1, -1, -1)
		assert.Error(t, err)
	})

	t.Run("Change game score of none-existent game should fail", func(t *testing.T) {
		testId1 := "test_id_2"
		match1, _ := NewMatch(&testId1, "one", "two")

		testId2 := "test_id_2"
		match2, _ := NewMatch(&testId2, "one", "two")

		round := NewRound(1)
		round.Matches = []Match{match1, match2}

		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Rounds:   []Round{round},
			Status:   SeasonStatusInProgress,
		}

		_, err := season.ChangeMatchScore("test_id_3", -1, -1)
		assert.Error(t, err)
	})

	t.Run("Change game score of not in progress season should fail", func(t *testing.T) {
		testId1 := "test_id_2"
		match1, _ := NewMatch(&testId1, "one", "two")

		testId2 := "test_id_2"
		match2, _ := NewMatch(&testId2, "one", "two")

		round := NewRound(1)
		round.Matches = []Match{match1, match2}

		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Rounds:   []Round{round},
			Status:   SeasonStatusPending,
		}

		_, err := season.ChangeMatchScore(testId1, -1, -1)
		assert.Error(t, err)
	})

	t.Run("Change game score of team playing against 'Bye' should fail", func(t *testing.T) {
		testId1 := "test_id_1"
		match1, _ := NewMatch(&testId1, "one", "bye")
		match1.Status = MatchStatusInProgress

		testId2 := "test_id_2"
		match2, _ := NewMatch(&testId2, "one", "two")
		match2.Status = MatchStatusInProgress

		round := NewRound(1)
		round.Matches = []Match{match1, match2}

		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Rounds:   []Round{round},
			Status:   SeasonStatusInProgress,
		}

		_, err := season.ChangeMatchScore(testId1, 1, 1)
		assert.Error(t, err)
	})

	t.Run("Change game score outside current round should fail", func(t *testing.T) {
		currentMatchID := "current-match"
		currentMatch, _ := NewMatch(&currentMatchID, "alpha", "beta")
		currentMatch.Status = MatchStatusInProgress

		testId1 := "next-round-match"
		match1, _ := NewMatch(&testId1, "one", "two")
		match1.Status = MatchStatusScheduled

		round := NewRound(1)
		round.Matches = []Match{currentMatch}

		nextRound := NewRound(2)
		nextRound.Matches = []Match{match1}

		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Rounds:   []Round{round, nextRound},
			Status:   SeasonStatusInProgress,
		}

		_, err := season.ChangeMatchScore(testId1, 1, 1)
		assert.Error(t, err)
	})

	t.Run("Change game score should pass", func(t *testing.T) {
		testId1 := "test_id_2"
		match1, _ := NewMatch(&testId1, "one", "two")
		match1.Status = MatchStatusInProgress

		testId2 := "test_id_2"
		match2, _ := NewMatch(&testId2, "one", "two")
		match2.Status = MatchStatusInProgress

		round := NewRound(1)
		round.Matches = []Match{match1, match2}

		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Rounds:   []Round{round},
			Status:   SeasonStatusInProgress,
		}

		changedSeason, _ := season.ChangeMatchScore(testId1, 2, 2)

		var findMatch = func(season Season, matchID string) *Match {
			for _, round := range season.Rounds {
				for _, match := range round.Matches {
					if matchID == testId1 {
						return &match
					}
				}
			}
			return nil
		}

		previousMatch := findMatch(season, testId1)
		currentMatch := findMatch(*changedSeason, testId1)

		assert.NotNil(t, currentMatch)
		assert.Equal(t, 2, currentMatch.HomeTeamScore)
		assert.Equal(t, 2, currentMatch.AwayTeamScore)
		assert.Equal(t, 0, previousMatch.HomeTeamScore)
		assert.Equal(t, 0, previousMatch.AwayTeamScore)
	})

	t.Run("Complete current round should advance next round", func(t *testing.T) {
		match1ID := "match-1"
		match1, _ := NewMatch(&match1ID, "one", "two")
		match1.Status = MatchStatusInProgress

		match2ID := "match-2"
		match2, _ := NewMatch(&match2ID, "three", "bye")
		match2.Status = MatchStatusInProgress

		match3ID := "match-3"
		match3, _ := NewMatch(&match3ID, "four", "five")

		round1 := NewRound(1)
		round1.Matches = []Match{match1, match2}

		round2 := NewRound(2)
		round2.Matches = []Match{match3}

		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Rounds:   []Round{round1, round2},
			Status:   SeasonStatusInProgress,
		}

		updatedSeason, err := season.CompleteCurrentRound()

		assert.NoError(t, err)
		assert.Equal(t, MatchStatusFinished, updatedSeason.Rounds[0].Matches[0].Status)
		assert.Equal(t, MatchStatusFinished, updatedSeason.Rounds[0].Matches[1].Status)
		assert.Equal(t, MatchStatusInProgress, updatedSeason.Rounds[1].Matches[0].Status)
		assert.Equal(t, SeasonStatusInProgress, updatedSeason.Status)
	})

	t.Run("Complete final round should finish season", func(t *testing.T) {
		match1ID := "match-1"
		match1, _ := NewMatch(&match1ID, "one", "two")
		match1.Status = MatchStatusInProgress

		round1 := NewRound(1)
		round1.Matches = []Match{match1}

		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Rounds:   []Round{round1},
			Status:   SeasonStatusInProgress,
		}

		updatedSeason, err := season.CompleteCurrentRound()

		assert.NoError(t, err)
		assert.Equal(t, MatchStatusFinished, updatedSeason.Rounds[0].Matches[0].Status)
		assert.Equal(t, SeasonStatusFinished, updatedSeason.Status)
	})

	t.Run("Configure playoff rules should allow finished regular season", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusFinished,
			Phase:    SeasonPhaseRegularSeason,
		}

		updatedSeason, err := season.ConfigurePlayoffRules(PlayoffRules{
			QualificationType: "top_n",
			QualifierCount:    4,
			Rounds: []PlayoffRoundRule{
				{Name: "semifinal", Legs: 2, TiedAggregateResolution: "penalties"},
				{Name: "final", Legs: 1, TiedAggregateResolution: "penalties"},
			},
		})

		assert.NoError(t, err)
		assert.NotNil(t, updatedSeason)
		assert.NotNil(t, updatedSeason.PlayoffRules)
		assert.Equal(t, 4, updatedSeason.PlayoffRules.QualifierCount)
	})

	t.Run("Configure playoff rules should fail once playoffs started", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusInProgress,
			Phase:    SeasonPhasePlayoffs,
			PlayoffBracket: &PlayoffBracket{
				Rounds: []PlayoffBracketRound{
					{
						Name:  "semifinal",
						Order: 1,
						Ties: []PlayoffTie{
							{
								ID:     "tie-1",
								Status: "in_progress",
								Matches: []Match{
									{ID: "match-1", PlayoffTieID: "tie-1", MatchOrder: 1, HomeTeamID: "team-1", AwayTeamID: "team-4", Status: MatchStatusFinished, HomeTeamScore: 1, AwayTeamScore: 0},
								},
							},
						},
					},
				},
			},
		}

		_, err := season.ConfigurePlayoffRules(PlayoffRules{
			QualificationType: "top_n",
			QualifierCount:    4,
			Rounds: []PlayoffRoundRule{
				{Name: "semifinal", Legs: 2, TiedAggregateResolution: "penalties"},
			},
		})

		assert.Error(t, err)
	})

	t.Run("Configure playoff rules should allow bracket reset before any playoff match is played", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusInProgress,
			Phase:    SeasonPhasePlayoffs,
			PlayoffRules: &PlayoffRules{
				QualificationType: "top_n",
				QualifierCount:    4,
				Rounds: []PlayoffRoundRule{
					{Name: "semifinal", Legs: 1, TiedAggregateResolution: "higher_seed_advances"},
				},
			},
			PlayoffBracket: &PlayoffBracket{
				Rounds: []PlayoffBracketRound{
					{
						Name:  "semifinal",
						Order: 1,
						Ties: []PlayoffTie{
							{
								ID:         "tie-1",
								RoundName:  "semifinal",
								RoundOrder: 1,
								SlotOrder:  1,
								Status:     "ready",
								Matches: []Match{
									{ID: "match-1", PlayoffTieID: "tie-1", MatchOrder: 1, HomeTeamID: "team-1", AwayTeamID: "team-4", Status: MatchStatusScheduled},
								},
							},
						},
					},
				},
			},
		}

		updatedSeason, err := season.ConfigurePlayoffRules(PlayoffRules{
			QualificationType: "top_n",
			QualifierCount:    2,
			Rounds: []PlayoffRoundRule{
				{Name: "final", Legs: 1, TiedAggregateResolution: "clear_winner_required"},
			},
		})

		assert.NoError(t, err)
		assert.NotNil(t, updatedSeason)
		assert.Nil(t, updatedSeason.PlayoffBracket)
		assert.Equal(t, SeasonPhaseRegularSeason, updatedSeason.Phase)
		assert.Equal(t, SeasonStatusFinished, updatedSeason.Status)
		assert.Equal(t, 2, updatedSeason.PlayoffRules.QualifierCount)
	})

	t.Run("Configure playoff rules should fail after a playoff match has been played", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusInProgress,
			Phase:    SeasonPhasePlayoffs,
			PlayoffBracket: &PlayoffBracket{
				Rounds: []PlayoffBracketRound{
					{
						Name:  "semifinal",
						Order: 1,
						Ties: []PlayoffTie{
							{
								ID:         "tie-1",
								RoundName:  "semifinal",
								RoundOrder: 1,
								SlotOrder:  1,
								Status:     "in_progress",
								Matches: []Match{
									{ID: "match-1", PlayoffTieID: "tie-1", MatchOrder: 1, HomeTeamID: "team-1", AwayTeamID: "team-4", Status: MatchStatusFinished, HomeTeamScore: 1, AwayTeamScore: 0},
								},
							},
						},
					},
				},
			},
		}

		_, err := season.ConfigurePlayoffRules(PlayoffRules{
			QualificationType: "top_n",
			QualifierCount:    4,
			Rounds: []PlayoffRoundRule{
				{Name: "semifinal", Legs: 1, TiedAggregateResolution: "higher_seed_advances"},
			},
		})

		assert.Error(t, err)
	})

	t.Run("Generate playoff bracket should move season into playoffs", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusFinished,
			Phase:    SeasonPhaseRegularSeason,
			PlayoffRules: &PlayoffRules{
				QualificationType: "top_n",
				QualifierCount:    4,
				Rounds: []PlayoffRoundRule{
					{Name: "semifinal", Legs: 2, HigherSeedHostsSecondLeg: true, TiedAggregateResolution: "penalties"},
					{Name: "final", Legs: 1, TiedAggregateResolution: "penalties"},
				},
			},
		}

		updatedSeason, err := season.GeneratePlayoffBracket([]PlayoffQualifiedTeam{
			{TeamID: "team-1", Seed: 1},
			{TeamID: "team-2", Seed: 2},
			{TeamID: "team-3", Seed: 3},
			{TeamID: "team-4", Seed: 4},
		})

		assert.NoError(t, err)
		assert.NotNil(t, updatedSeason.PlayoffBracket)
		assert.Equal(t, SeasonPhasePlayoffs, updatedSeason.Phase)
		assert.Equal(t, SeasonStatusInProgress, updatedSeason.Status)
		assert.Len(t, updatedSeason.PlayoffBracket.Rounds, 2)
		assert.Len(t, updatedSeason.PlayoffBracket.Rounds[0].Ties[0].Matches, 2)
	})

	t.Run("Generate playoff bracket should replace invalid unstarted bracket", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusFinished,
			Phase:    SeasonPhaseRegularSeason,
			PlayoffRules: &PlayoffRules{
				QualificationType: "top_n",
				QualifierCount:    4,
				Rounds: []PlayoffRoundRule{
					{Name: "semifinal", Legs: 2, HigherSeedHostsSecondLeg: true, TiedAggregateResolution: "higher_seed_advances"},
					{Name: "final", Legs: 1, TiedAggregateResolution: "clear_winner_required"},
				},
			},
			PlayoffBracket: &PlayoffBracket{
				Rounds: []PlayoffBracketRound{
					{
						Name:  "semifinal",
						Order: 1,
						Ties: []PlayoffTie{
							{
								ID:         "old-tie-1",
								RoundName:  "semifinal",
								RoundOrder: 1,
								SlotOrder:  1,
								HomeSeed:   1,
								AwaySeed:   4,
								HomeTeamID: "team-1",
								AwayTeamID: "team-4",
								Status:     "pending",
								Matches:    []Match{},
							},
						},
					},
				},
			},
		}

		updatedSeason, err := season.GeneratePlayoffBracket([]PlayoffQualifiedTeam{
			{TeamID: "team-1", Seed: 1},
			{TeamID: "team-2", Seed: 2},
			{TeamID: "team-3", Seed: 3},
			{TeamID: "team-4", Seed: 4},
		})

		assert.NoError(t, err)
		assert.NotNil(t, updatedSeason.PlayoffBracket)
		assert.Len(t, updatedSeason.PlayoffBracket.Rounds[0].Ties, 2)
		assert.Len(t, updatedSeason.PlayoffBracket.Rounds[0].Ties[0].Matches, 2)
		assert.Equal(t, "ready", updatedSeason.PlayoffBracket.Rounds[0].Ties[0].Status)
	})

	t.Run("Record playoff match score should finish match", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusInProgress,
			Phase:    SeasonPhasePlayoffs,
			PlayoffBracket: &PlayoffBracket{
				Rounds: []PlayoffBracketRound{
					{
						Name:  "semifinal",
						Order: 1,
						Ties: []PlayoffTie{
							{
								ID:         "tie-1",
								RoundName:  "semifinal",
								RoundOrder: 1,
								SlotOrder:  1,
								HomeTeamID: "team-1",
								AwayTeamID: "team-4",
								Status:     "pending",
								Matches: []Match{
									{ID: "leg-1", PlayoffTieID: "tie-1", MatchOrder: 1, HomeTeamID: "team-4", AwayTeamID: "team-1", Status: MatchStatusScheduled},
									{ID: "leg-2", PlayoffTieID: "tie-1", MatchOrder: 2, HomeTeamID: "team-1", AwayTeamID: "team-4", Status: MatchStatusScheduled},
								},
							},
						},
					},
				},
			},
		}

		updatedSeason, err := season.RecordPlayoffMatchScore("tie-1", "leg-1", 2, 1)

		assert.NoError(t, err)
		assert.Equal(t, 2, updatedSeason.PlayoffBracket.Rounds[0].Ties[0].Matches[0].HomeTeamScore)
		assert.Equal(t, MatchStatusFinished, updatedSeason.PlayoffBracket.Rounds[0].Ties[0].Matches[0].Status)
		assert.Equal(t, "in_progress", updatedSeason.PlayoffBracket.Rounds[0].Ties[0].Status)
	})

	t.Run("Record playoff final match should set champion when aggregate winner exists", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusInProgress,
			Phase:    SeasonPhasePlayoffs,
			PlayoffRules: &PlayoffRules{
				QualificationType: "top_n",
				QualifierCount:    2,
				Rounds: []PlayoffRoundRule{
					{Name: "final", Legs: 1, TiedAggregateResolution: "penalties"},
				},
			},
			PlayoffBracket: &PlayoffBracket{
				Rounds: []PlayoffBracketRound{
					{
						Name:  "final",
						Order: 1,
						Ties: []PlayoffTie{
							{
								ID:         "tie-final",
								RoundName:  "final",
								RoundOrder: 1,
								SlotOrder:  1,
								HomeSeed:   1,
								AwaySeed:   2,
								HomeTeamID: "team-1",
								AwayTeamID: "team-2",
								Status:     "ready",
								Matches: []Match{
									{ID: "leg-final", PlayoffTieID: "tie-final", MatchOrder: 1, HomeTeamID: "team-1", AwayTeamID: "team-2", Status: MatchStatusScheduled},
								},
							},
						},
					},
				},
			},
		}

		updatedSeason, err := season.RecordPlayoffMatchScore("tie-final", "leg-final", 3, 1)

		assert.NoError(t, err)
		assert.NotNil(t, updatedSeason.ChampionTeamID)
		assert.Equal(t, "team-1", *updatedSeason.ChampionTeamID)
		assert.Equal(t, SeasonPhaseCompleted, updatedSeason.Phase)
		assert.Equal(t, SeasonStatusFinished, updatedSeason.Status)
	})

	t.Run("Record playoff semifinal should advance winner into final", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusInProgress,
			Phase:    SeasonPhasePlayoffs,
			PlayoffRules: &PlayoffRules{
				QualificationType: "top_n",
				QualifierCount:    4,
				Rounds: []PlayoffRoundRule{
					{Name: "semifinal", Legs: 1, TiedAggregateResolution: "penalties"},
					{Name: "final", Legs: 1, TiedAggregateResolution: "penalties"},
				},
			},
			PlayoffBracket: &PlayoffBracket{
				Rounds: []PlayoffBracketRound{
					{
						Name:  "semifinal",
						Order: 1,
						Ties: []PlayoffTie{
							{
								ID:         "semi-1",
								RoundName:  "semifinal",
								RoundOrder: 1,
								SlotOrder:  1,
								HomeSeed:   1,
								AwaySeed:   4,
								HomeTeamID: "team-1",
								AwayTeamID: "team-4",
								Status:     "ready",
								Matches: []Match{
									{ID: "semi-leg-1", PlayoffTieID: "semi-1", MatchOrder: 1, HomeTeamID: "team-1", AwayTeamID: "team-4", Status: MatchStatusScheduled},
								},
							},
							{
								ID:         "semi-2",
								RoundName:  "semifinal",
								RoundOrder: 1,
								SlotOrder:  2,
								HomeSeed:   2,
								AwaySeed:   3,
								HomeTeamID: "team-2",
								AwayTeamID: "team-3",
								Status:     "ready",
								Matches: []Match{
									{ID: "semi-leg-2", PlayoffTieID: "semi-2", MatchOrder: 1, HomeTeamID: "team-2", AwayTeamID: "team-3", Status: MatchStatusScheduled},
								},
							},
						},
					},
					{
						Name:  "final",
						Order: 2,
						Ties: []PlayoffTie{
							{
								ID:         "final-1",
								RoundName:  "final",
								RoundOrder: 2,
								SlotOrder:  1,
								Status:     "pending",
								Matches:    []Match{},
							},
						},
					},
				},
			},
		}

		updatedSeason, err := season.RecordPlayoffMatchScore("semi-1", "semi-leg-1", 2, 0)

		assert.NoError(t, err)
		assert.NotNil(t, updatedSeason.PlayoffBracket.Rounds[0].Ties[0].WinnerTeamID)
		assert.Equal(t, "team-1", *updatedSeason.PlayoffBracket.Rounds[0].Ties[0].WinnerTeamID)
		assert.Equal(t, "team-1", updatedSeason.PlayoffBracket.Rounds[1].Ties[0].HomeTeamID)
		assert.Empty(t, updatedSeason.PlayoffBracket.Rounds[1].Ties[0].Matches)
	})

	t.Run("Record second playoff semifinal should advance winner into final away slot and create final matches", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusInProgress,
			Phase:    SeasonPhasePlayoffs,
			PlayoffRules: &PlayoffRules{
				QualificationType: "top_n",
				QualifierCount:    4,
				Rounds: []PlayoffRoundRule{
					{Name: "semifinal", Legs: 1, TiedAggregateResolution: "higher_seed_advances"},
					{Name: "final", Legs: 1, TiedAggregateResolution: "clear_winner_required"},
				},
			},
			PlayoffBracket: &PlayoffBracket{
				Rounds: []PlayoffBracketRound{
					{
						Name:  "semifinal",
						Order: 1,
						Ties: []PlayoffTie{
							{
								ID:         "semi-1",
								RoundName:  "semifinal",
								RoundOrder: 1,
								SlotOrder:  1,
								HomeSeed:   1,
								AwaySeed:   4,
								HomeTeamID: "team-1",
								AwayTeamID: "team-4",
								Status:     "finished",
								WinnerTeamID: func() *string {
									v := "team-1"
									return &v
								}(),
								Matches: []Match{
									{ID: "semi-leg-1", PlayoffTieID: "semi-1", MatchOrder: 1, HomeTeamID: "team-1", AwayTeamID: "team-4", Status: MatchStatusFinished, HomeTeamScore: 2, AwayTeamScore: 0},
								},
							},
							{
								ID:         "semi-2",
								RoundName:  "semifinal",
								RoundOrder: 1,
								SlotOrder:  2,
								HomeSeed:   2,
								AwaySeed:   3,
								HomeTeamID: "team-2",
								AwayTeamID: "team-3",
								Status:     "ready",
								Matches: []Match{
									{ID: "semi-leg-2", PlayoffTieID: "semi-2", MatchOrder: 1, HomeTeamID: "team-2", AwayTeamID: "team-3", Status: MatchStatusScheduled},
								},
							},
						},
					},
					{
						Name:  "final",
						Order: 2,
						Ties: []PlayoffTie{
							{
								ID:         "final-1",
								RoundName:  "final",
								RoundOrder: 2,
								SlotOrder:  1,
								HomeTeamID: "team-1",
								HomeSeed:   1,
								Status:     "pending",
								Matches:    []Match{},
							},
						},
					},
				},
			},
		}

		updatedSeason, err := season.RecordPlayoffMatchScore("semi-2", "semi-leg-2", 1, 0)

		assert.NoError(t, err)
		assert.NotNil(t, updatedSeason.PlayoffBracket.Rounds[0].Ties[1].WinnerTeamID)
		assert.Equal(t, "team-2", *updatedSeason.PlayoffBracket.Rounds[0].Ties[1].WinnerTeamID)
		assert.Equal(t, "team-2", updatedSeason.PlayoffBracket.Rounds[1].Ties[0].AwayTeamID)
		assert.Equal(t, 2, updatedSeason.PlayoffBracket.Rounds[1].Ties[0].AwaySeed)
		assert.Equal(t, "ready", updatedSeason.PlayoffBracket.Rounds[1].Ties[0].Status)
		assert.Len(t, updatedSeason.PlayoffBracket.Rounds[1].Ties[0].Matches, 1)
	})

	t.Run("Record tied playoff semifinal should advance higher seed", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusInProgress,
			Phase:    SeasonPhasePlayoffs,
			PlayoffRules: &PlayoffRules{
				QualificationType: "top_n",
				QualifierCount:    4,
				Rounds: []PlayoffRoundRule{
					{Name: "semifinal", Legs: 2, HigherSeedHostsSecondLeg: true, TiedAggregateResolution: "higher_seed_advances"},
					{Name: "final", Legs: 1, TiedAggregateResolution: "clear_winner_required"},
				},
			},
			PlayoffBracket: &PlayoffBracket{
				Rounds: []PlayoffBracketRound{
					{
						Name:  "semifinal",
						Order: 1,
						Ties: []PlayoffTie{
							{
								ID:         "semi-1",
								RoundName:  "semifinal",
								RoundOrder: 1,
								SlotOrder:  1,
								HomeSeed:   1,
								AwaySeed:   4,
								HomeTeamID: "team-1",
								AwayTeamID: "team-4",
								Status:     "in_progress",
								Matches: []Match{
									{ID: "semi-leg-1", PlayoffTieID: "semi-1", MatchOrder: 1, HomeTeamID: "team-4", AwayTeamID: "team-1", Status: MatchStatusFinished, HomeTeamScore: 2, AwayTeamScore: 1},
									{ID: "semi-leg-2", PlayoffTieID: "semi-1", MatchOrder: 2, HomeTeamID: "team-1", AwayTeamID: "team-4", Status: MatchStatusScheduled},
								},
							},
							{
								ID:         "semi-2",
								RoundName:  "semifinal",
								RoundOrder: 1,
								SlotOrder:  2,
								HomeSeed:   2,
								AwaySeed:   3,
								HomeTeamID: "team-2",
								AwayTeamID: "team-3",
								Status:     "pending",
							},
						},
					},
					{
						Name:  "final",
						Order: 2,
						Ties: []PlayoffTie{
							{
								ID:         "final-1",
								RoundName:  "final",
								RoundOrder: 2,
								SlotOrder:  1,
								Status:     "pending",
							},
						},
					},
				},
			},
		}

		updatedSeason, err := season.RecordPlayoffMatchScore("semi-1", "semi-leg-2", 1, 0)

		assert.NoError(t, err)
		assert.NotNil(t, updatedSeason.PlayoffBracket.Rounds[0].Ties[0].WinnerTeamID)
		assert.Equal(t, "team-1", *updatedSeason.PlayoffBracket.Rounds[0].Ties[0].WinnerTeamID)
		assert.Equal(t, "finished", updatedSeason.PlayoffBracket.Rounds[0].Ties[0].Status)
		assert.Equal(t, "team-1", updatedSeason.PlayoffBracket.Rounds[1].Ties[0].HomeTeamID)
	})

	t.Run("Record tied playoff final should fail because final must have winner", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
			Status:   SeasonStatusInProgress,
			Phase:    SeasonPhasePlayoffs,
			PlayoffRules: &PlayoffRules{
				QualificationType: "top_n",
				QualifierCount:    2,
				Rounds: []PlayoffRoundRule{
					{Name: "final", Legs: 1, TiedAggregateResolution: "clear_winner_required"},
				},
			},
			PlayoffBracket: &PlayoffBracket{
				Rounds: []PlayoffBracketRound{
					{
						Name:  "final",
						Order: 1,
						Ties: []PlayoffTie{
							{
								ID:         "tie-final",
								RoundName:  "final",
								RoundOrder: 1,
								SlotOrder:  1,
								HomeSeed:   1,
								AwaySeed:   2,
								HomeTeamID: "team-1",
								AwayTeamID: "team-2",
								Status:     "ready",
								Matches: []Match{
									{ID: "leg-final", PlayoffTieID: "tie-final", MatchOrder: 1, HomeTeamID: "team-1", AwayTeamID: "team-2", Status: MatchStatusScheduled},
								},
							},
						},
					},
				},
			},
		}

		updatedSeason, err := season.RecordPlayoffMatchScore("tie-final", "leg-final", 1, 1)

		assert.Error(t, err)
		assert.Nil(t, updatedSeason)
	})

}
