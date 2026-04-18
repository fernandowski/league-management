package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSeason_ScheduleRounds(t *testing.T) {
	t.Run("Zero teams should error out", func(t *testing.T) {
		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusPending, "", 0, nil, nil, nil, nil)

		leagueId := "league-id-1"
		league := League{
			Id:             &leagueId,
			Name:           "test-league",
			OwnerId:        "test-owner-id",
			OrganizationId: "test-org-id",
			Memberships:    []LeagueMembership{},
		}

		err := season.ScheduleRounds(leagueId, league.Memberships)
		assert.Error(t, err)
	})

	t.Run("Season with even teams without rules", func(t *testing.T) {
		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusPending, "", 0, nil, nil, nil, nil)

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

		leagueId := "league-id-1"
		league := League{
			Id:             &leagueId,
			Name:           "test-league",
			OwnerId:        "test-owner-id",
			OrganizationId: "test-org-id",
			Memberships:    []LeagueMembership{member1, member2, member3, member4},
		}

		err := season.ScheduleRounds(leagueId, league.Memberships)
		firstRound, ok := season.FindRound(1)

		assert.NoError(t, err)
		assert.Equal(t, 3, season.RoundCount())
		assert.Equal(t, season.CurrentStatus(), SeasonStatusPlanned)
		assert.True(t, ok)
		assert.Equal(t, 1, firstRound.RoundNumber)
		assert.Len(t, firstRound.Matches, 2)
	})

	t.Run("Season with odd teams without rules", func(t *testing.T) {
		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusPending, "", 0, nil, nil, nil, nil)

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

		leagueId := "league-id-1"
		league := League{
			Id:             &leagueId,
			Name:           "test-league",
			OwnerId:        "test-owner-id",
			OrganizationId: "test-org-id",
			Memberships:    []LeagueMembership{member1, member2, member3},
		}

		err := season.ScheduleRounds(leagueId, league.Memberships)
		firstRound, ok := season.FindRound(1)

		assert.NoError(t, err)
		assert.Equal(t, 3, season.RoundCount())
		assert.Equal(t, season.CurrentStatus(), SeasonStatusPlanned)
		assert.True(t, ok)
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
		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusPending, "", 0, nil, nil, nil, nil)

		_, err := season.Start()
		assert.Error(t, err)
	})

	t.Run("Start season in_progress status should fail", func(t *testing.T) {
		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusInProgress, "", 0, nil, nil, nil, nil)

		_, err := season.Start()
		assert.Error(t, err)
	})

	t.Run("Start season paused status should fail", func(t *testing.T) {
		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusPaused, "", 0, nil, nil, nil, nil)

		_, err := season.Start()
		assert.Error(t, err)
	})

	t.Run("Start season finished status should fail", func(t *testing.T) {
		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusFinished, "", 0, nil, nil, nil, nil)

		_, err := season.Start()
		assert.Error(t, err)
	})

	t.Run("Start season undefined status should fail", func(t *testing.T) {
		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusUndefined, "", 0, nil, nil, nil, nil)

		_, err := season.Start()
		assert.Error(t, err)
	})

	t.Run("Start Season with empty rounds should fail", func(t *testing.T) {
		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusPending, "", 0, []Round{}, nil, nil, nil)

		_, err := season.Start()
		assert.Error(t, err)
	})

	t.Run("Season with even teams without rules", func(t *testing.T) {
		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusPending, "", 0, nil, nil, nil, nil)

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

		leagueId := "league-id-1"
		league := League{
			Id:             &leagueId,
			Name:           "test-league",
			OwnerId:        "test-owner-id",
			OrganizationId: "test-org-id",
			Memberships:    []LeagueMembership{member1, member2, member3, member4},
		}

		_ = season.ScheduleRounds(leagueId, league.Memberships)
		newSeason, _ := season.Start()

		assert.Equal(t, SeasonStatusInProgress, newSeason.CurrentStatus())
		assert.Equal(t, SeasonStatusPlanned, season.CurrentStatus())
		currentRound, ok := newSeason.CurrentRound()
		assert.True(t, ok)
		for _, match := range currentRound.Matches {
			assert.Equal(t, MatchStatusInProgress, match.Status)
		}
	})

	t.Run("Schedule rounds with different league should fail", func(t *testing.T) {
		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusPending, "", 0, nil, nil, nil, nil)

		err := season.ScheduleRounds("league-id-2", []LeagueMembership{
			{ID: "m1", TeamID: "team-a"},
			{ID: "m2", TeamID: "team-b"},
		})

		assert.Error(t, err)
	})

	t.Run("Change game score to negative number should fail", func(t *testing.T) {

		testId1 := "test_id_2"
		match1, _ := NewMatch(&testId1, "one", "two")

		testId2 := "test_id_2"
		match2, _ := NewMatch(&testId2, "one", "two")

		round := RehydrateRound(1, []Match{match1, match2})

		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusInProgress, "", 0, []Round{round}, nil, nil, nil)

		_, err := season.ChangeMatchScore(testId1, -1, -1)
		assert.Error(t, err)
	})

	t.Run("Change game score of none-existent game should fail", func(t *testing.T) {
		testId1 := "test_id_2"
		match1, _ := NewMatch(&testId1, "one", "two")

		testId2 := "test_id_2"
		match2, _ := NewMatch(&testId2, "one", "two")

		round := RehydrateRound(1, []Match{match1, match2})

		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusInProgress, "", 0, []Round{round}, nil, nil, nil)

		_, err := season.ChangeMatchScore("test_id_3", -1, -1)
		assert.Error(t, err)
	})

	t.Run("Change game score of not in progress season should fail", func(t *testing.T) {
		testId1 := "test_id_2"
		match1, _ := NewMatch(&testId1, "one", "two")

		testId2 := "test_id_2"
		match2, _ := NewMatch(&testId2, "one", "two")

		round := RehydrateRound(1, []Match{match1, match2})

		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusPending, "", 0, []Round{round}, nil, nil, nil)

		_, err := season.ChangeMatchScore(testId1, -1, -1)
		assert.Error(t, err)
	})

	t.Run("Change game score of team playing against 'Bye' should fail", func(t *testing.T) {
		testId1 := "test_id_1"
		match1 := RehydrateMatch(MatchState{
			ID:         testId1,
			HomeTeamID: "one",
			AwayTeamID: "bye",
			Status:     MatchStatusInProgress,
		})

		testId2 := "test_id_2"
		match2 := RehydrateMatch(MatchState{
			ID:         testId2,
			HomeTeamID: "one",
			AwayTeamID: "two",
			Status:     MatchStatusInProgress,
		})

		round := RehydrateRound(1, []Match{match1, match2})

		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusInProgress, "", 0, []Round{round}, nil, nil, nil)

		_, err := season.ChangeMatchScore(testId1, 1, 1)
		assert.Error(t, err)
	})

	t.Run("Change game score outside current round should fail", func(t *testing.T) {
		currentMatchID := "current-match"
		currentMatch := RehydrateMatch(MatchState{
			ID:         currentMatchID,
			HomeTeamID: "alpha",
			AwayTeamID: "beta",
			Status:     MatchStatusInProgress,
		})

		testId1 := "next-round-match"
		match1 := RehydrateMatch(MatchState{
			ID:         testId1,
			HomeTeamID: "one",
			AwayTeamID: "two",
			Status:     MatchStatusScheduled,
		})

		round := RehydrateRound(1, []Match{currentMatch})

		nextRound := RehydrateRound(2, []Match{match1})

		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusInProgress, "", 0, []Round{round, nextRound}, nil, nil, nil)

		_, err := season.ChangeMatchScore(testId1, 1, 1)
		assert.Error(t, err)
	})

	t.Run("Change game score should pass", func(t *testing.T) {
		testId1 := "test_id_2"
		match1 := RehydrateMatch(MatchState{
			ID:         testId1,
			HomeTeamID: "one",
			AwayTeamID: "two",
			Status:     MatchStatusInProgress,
		})

		testId2 := "test_id_2"
		match2 := RehydrateMatch(MatchState{
			ID:         testId2,
			HomeTeamID: "one",
			AwayTeamID: "two",
			Status:     MatchStatusInProgress,
		})

		round := RehydrateRound(1, []Match{match1, match2})

		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusInProgress, "", 0, []Round{round}, nil, nil, nil)

		changedSeason, _ := season.ChangeMatchScore(testId1, 2, 2)

		previousMatch, _ := season.FindMatchSnapshot(testId1)
		currentMatch, _ := changedSeason.FindMatchSnapshot(testId1)

		assert.NotNil(t, currentMatch)
		assert.Equal(t, 2, currentMatch.HomeTeamScore)
		assert.Equal(t, 2, currentMatch.AwayTeamScore)
		assert.Equal(t, 0, previousMatch.HomeTeamScore)
		assert.Equal(t, 0, previousMatch.AwayTeamScore)
	})

	t.Run("Complete current round should advance next round", func(t *testing.T) {
		match1ID := "match-1"
		match1 := RehydrateMatch(MatchState{
			ID:         match1ID,
			HomeTeamID: "one",
			AwayTeamID: "two",
			Status:     MatchStatusInProgress,
		})

		match2ID := "match-2"
		match2 := RehydrateMatch(MatchState{
			ID:         match2ID,
			HomeTeamID: "three",
			AwayTeamID: "bye",
			Status:     MatchStatusInProgress,
		})

		match3ID := "match-3"
		match3, _ := NewMatch(&match3ID, "four", "five")

		round1 := RehydrateRound(1, []Match{match1, match2})

		round2 := RehydrateRound(2, []Match{match3})

		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusInProgress, "", 0, []Round{round1, round2}, nil, nil, nil)

		updatedSeason, err := season.CompleteCurrentRound()

		assert.NoError(t, err)
		currentRound, ok := updatedSeason.CurrentRound()
		assert.True(t, ok)
		roundOne, ok := updatedSeason.FindRound(1)
		assert.True(t, ok)
		assert.Equal(t, MatchStatusFinished, roundOne.Matches[0].Status)
		assert.Equal(t, MatchStatusFinished, roundOne.Matches[1].Status)
		assert.Equal(t, MatchStatusInProgress, currentRound.Matches[0].Status)
		assert.Equal(t, SeasonStatusInProgress, updatedSeason.CurrentStatus())
	})

	t.Run("Complete final round should finish season", func(t *testing.T) {
		match1ID := "match-1"
		match1 := RehydrateMatch(MatchState{
			ID:         match1ID,
			HomeTeamID: "one",
			AwayTeamID: "two",
			Status:     MatchStatusInProgress,
		})

		round1 := RehydrateRound(1, []Match{match1})

		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusInProgress, "", 0, []Round{round1}, nil, nil, nil)

		updatedSeason, err := season.CompleteCurrentRound()

		assert.NoError(t, err)
		matchSnapshot, ok := updatedSeason.FindMatchSnapshot(match1ID)
		assert.True(t, ok)
		assert.Equal(t, MatchStatusFinished, matchSnapshot.Status)
		assert.Equal(t, SeasonStatusFinished, updatedSeason.CurrentStatus())
	})

	t.Run("Configure playoff rules should allow finished regular season", func(t *testing.T) {
		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusFinished, SeasonPhaseRegularSeason, 0, nil, nil, nil, nil)

		updatedSeason, err := season.ConfigurePlayoffRules(PlayoffRulesSnapshot{
			QualificationType: "top_n",
			QualifierCount:    4,
			Rounds: []PlayoffRoundRuleSnapshot{
				{Name: "semifinal", Legs: 2, TiedAggregateResolution: "penalties"},
				{Name: "final", Legs: 1, TiedAggregateResolution: "penalties"},
			},
		})

		assert.NoError(t, err)
		assert.NotNil(t, updatedSeason)
		qualifierCount, configured := updatedSeason.PlayoffQualifierCount()
		assert.True(t, configured)
		assert.Equal(t, 4, qualifierCount)
	})

	t.Run("Configure playoff rules should fail once playoffs started", func(t *testing.T) {
		season := rehydratedSeasonForTest(
			"id-1",
			"league-id-1",
			"Test League",
			SeasonStatusInProgress,
			SeasonPhasePlayoffs,
			0,
			nil,
			nil,
			RehydratePlayoffBracket([]PlayoffBracketRoundSnapshot{
				{
					Name:  "semifinal",
					Order: 1,
					Ties: []PlayoffTieSnapshot{
						{
							ID:     "tie-1",
							Status: "in_progress",
							Matches: []MatchSnapshot{
								RehydrateMatch(MatchState{
									ID:            "match-1",
									PlayoffTieID:  "tie-1",
									MatchOrder:    1,
									HomeTeamID:    "team-1",
									AwayTeamID:    "team-4",
									Status:        MatchStatusFinished,
									HomeTeamScore: 1,
									AwayTeamScore: 0,
								}).Snapshot(),
							},
						},
					},
				},
			}),
			nil,
		)

		_, err := season.ConfigurePlayoffRules(PlayoffRulesSnapshot{
			QualificationType: "top_n",
			QualifierCount:    4,
			Rounds: []PlayoffRoundRuleSnapshot{
				{Name: "semifinal", Legs: 2, TiedAggregateResolution: "penalties"},
			},
		})

		assert.Error(t, err)
	})

	t.Run("Configure playoff rules should allow bracket reset before any playoff match is played", func(t *testing.T) {
		season := rehydratedSeasonForTest(
			"id-1",
			"league-id-1",
			"Test League",
			SeasonStatusInProgress,
			SeasonPhasePlayoffs,
			0,
			nil,
			&PlayoffRulesSnapshot{
				QualificationType: "top_n",
				QualifierCount:    4,
				Rounds: []PlayoffRoundRuleSnapshot{
					{Name: "semifinal", Legs: 1, TiedAggregateResolution: "higher_seed_advances"},
				},
			},
			RehydratePlayoffBracket([]PlayoffBracketRoundSnapshot{
				{
					Name:  "semifinal",
					Order: 1,
					Ties: []PlayoffTieSnapshot{
						{
							ID:         "tie-1",
							RoundName:  "semifinal",
							RoundOrder: 1,
							SlotOrder:  1,
							Status:     "ready",
							Matches: []MatchSnapshot{
								RehydrateMatch(MatchState{ID: "match-1", PlayoffTieID: "tie-1", MatchOrder: 1, HomeTeamID: "team-1", AwayTeamID: "team-4", Status: MatchStatusScheduled}).Snapshot(),
							},
						},
					},
				},
			}),
			nil,
		)

		updatedSeason, err := season.ConfigurePlayoffRules(PlayoffRulesSnapshot{
			QualificationType: "top_n",
			QualifierCount:    2,
			Rounds: []PlayoffRoundRuleSnapshot{
				{Name: "final", Legs: 1, TiedAggregateResolution: "clear_winner_required"},
			},
		})

		assert.NoError(t, err)
		assert.NotNil(t, updatedSeason)
		assert.False(t, updatedSeason.HasPlayoffBracket())
		assert.Equal(t, SeasonPhaseRegularSeason, updatedSeason.CurrentPhase())
		assert.Equal(t, SeasonStatusFinished, updatedSeason.CurrentStatus())
		qualifierCount, configured := updatedSeason.PlayoffQualifierCount()
		assert.True(t, configured)
		assert.Equal(t, 2, qualifierCount)
	})

	t.Run("Configure playoff rules should fail after a playoff match has been played", func(t *testing.T) {
		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusInProgress, SeasonPhasePlayoffs, 0, nil, nil, RehydratePlayoffBracket([]PlayoffBracketRoundSnapshot{
			{
				Name:  "semifinal",
				Order: 1,
				Ties: []PlayoffTieSnapshot{
					{
						ID:         "tie-1",
						RoundName:  "semifinal",
						RoundOrder: 1,
						SlotOrder:  1,
						Status:     "in_progress",
						Matches: []MatchSnapshot{
							RehydrateMatch(MatchState{ID: "match-1", PlayoffTieID: "tie-1", MatchOrder: 1, HomeTeamID: "team-1", AwayTeamID: "team-4", Status: MatchStatusFinished, HomeTeamScore: 1, AwayTeamScore: 0}).Snapshot(),
						},
					},
				},
			},
		}), nil)

		_, err := season.ConfigurePlayoffRules(PlayoffRulesSnapshot{
			QualificationType: "top_n",
			QualifierCount:    4,
			Rounds: []PlayoffRoundRuleSnapshot{
				{Name: "semifinal", Legs: 1, TiedAggregateResolution: "higher_seed_advances"},
			},
		})

		assert.Error(t, err)
	})

	t.Run("Generate playoff bracket should move season into playoffs", func(t *testing.T) {
		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusFinished, SeasonPhaseRegularSeason, 0, nil, &PlayoffRulesSnapshot{
			QualificationType: "top_n",
			QualifierCount:    4,
			Rounds: []PlayoffRoundRuleSnapshot{
				{Name: "semifinal", Legs: 2, HigherSeedHostsSecondLeg: true, TiedAggregateResolution: "penalties"},
				{Name: "final", Legs: 1, TiedAggregateResolution: "penalties"},
			},
		}, nil, nil)

		updatedSeason, err := season.GeneratePlayoffBracket([]PlayoffQualifiedTeam{
			{TeamID: "team-1", Seed: 1},
			{TeamID: "team-2", Seed: 2},
			{TeamID: "team-3", Seed: 3},
			{TeamID: "team-4", Seed: 4},
		})

		assert.NoError(t, err)
		assert.Equal(t, SeasonPhasePlayoffs, updatedSeason.CurrentPhase())
		assert.Equal(t, SeasonStatusInProgress, updatedSeason.CurrentStatus())
		firstRound, ok := updatedSeason.FindPlayoffRound(1)
		assert.True(t, ok)
		assert.Len(t, updatedSeason.PlayoffBracketRounds(), 2)
		assert.Len(t, firstRound.Ties[0].Matches, 2)
	})

	t.Run("Generate playoff bracket should replace invalid unstarted bracket", func(t *testing.T) {
		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusFinished, SeasonPhaseRegularSeason, 0, nil, &PlayoffRulesSnapshot{
			QualificationType: "top_n",
			QualifierCount:    4,
			Rounds: []PlayoffRoundRuleSnapshot{
				{Name: "semifinal", Legs: 2, HigherSeedHostsSecondLeg: true, TiedAggregateResolution: "higher_seed_advances"},
				{Name: "final", Legs: 1, TiedAggregateResolution: "clear_winner_required"},
			},
		}, RehydratePlayoffBracket([]PlayoffBracketRoundSnapshot{
			{
				Name:  "semifinal",
				Order: 1,
				Ties: []PlayoffTieSnapshot{
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
						Matches:    []MatchSnapshot{},
					},
				},
			},
		}), nil)

		updatedSeason, err := season.GeneratePlayoffBracket([]PlayoffQualifiedTeam{
			{TeamID: "team-1", Seed: 1},
			{TeamID: "team-2", Seed: 2},
			{TeamID: "team-3", Seed: 3},
			{TeamID: "team-4", Seed: 4},
		})

		assert.NoError(t, err)
		firstRound, ok := updatedSeason.FindPlayoffRound(1)
		assert.True(t, ok)
		assert.Len(t, firstRound.Ties, 2)
		assert.Len(t, firstRound.Ties[0].Matches, 2)
		assert.Equal(t, "ready", firstRound.Ties[0].Status)
	})

	t.Run("Record playoff match score should finish match", func(t *testing.T) {
		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusInProgress, SeasonPhasePlayoffs, 0, nil, nil, RehydratePlayoffBracket([]PlayoffBracketRoundSnapshot{
			{
				Name:  "semifinal",
				Order: 1,
				Ties: []PlayoffTieSnapshot{
					{
						ID:         "tie-1",
						RoundName:  "semifinal",
						RoundOrder: 1,
						SlotOrder:  1,
						HomeTeamID: "team-1",
						AwayTeamID: "team-4",
						Status:     "pending",
						Matches: []MatchSnapshot{
							RehydrateMatch(MatchState{ID: "leg-1", PlayoffTieID: "tie-1", MatchOrder: 1, HomeTeamID: "team-4", AwayTeamID: "team-1", Status: MatchStatusScheduled}).Snapshot(),
							RehydrateMatch(MatchState{ID: "leg-2", PlayoffTieID: "tie-1", MatchOrder: 2, HomeTeamID: "team-1", AwayTeamID: "team-4", Status: MatchStatusScheduled}).Snapshot(),
						},
					},
				},
			},
		}), nil)

		updatedSeason, err := season.RecordPlayoffMatchScore("tie-1", "leg-1", 2, 1)

		assert.NoError(t, err)
		tie, ok := updatedSeason.FindPlayoffTie("tie-1")
		assert.True(t, ok)
		assert.Equal(t, 2, tie.Matches[0].HomeTeamScore)
		assert.Equal(t, MatchStatusFinished, tie.Matches[0].Status)
		assert.Equal(t, "in_progress", tie.Status)
	})

	t.Run("Record playoff final match should set champion when aggregate winner exists", func(t *testing.T) {
		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusInProgress, SeasonPhasePlayoffs, 0, nil, &PlayoffRulesSnapshot{
			QualificationType: "top_n",
			QualifierCount:    2,
			Rounds: []PlayoffRoundRuleSnapshot{
				{Name: "final", Legs: 1, TiedAggregateResolution: "penalties"},
			},
		}, RehydratePlayoffBracket([]PlayoffBracketRoundSnapshot{
			{
				Name:  "final",
				Order: 1,
				Ties: []PlayoffTieSnapshot{
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
						Matches: []MatchSnapshot{
							RehydrateMatch(MatchState{ID: "leg-final", PlayoffTieID: "tie-final", MatchOrder: 1, HomeTeamID: "team-1", AwayTeamID: "team-2", Status: MatchStatusScheduled}).Snapshot(),
						},
					},
				},
			},
		}), nil)

		updatedSeason, err := season.RecordPlayoffMatchScore("tie-final", "leg-final", 3, 1)

		assert.NoError(t, err)
		assert.NotNil(t, updatedSeason.ChampionTeam())
		assert.Equal(t, "team-1", *updatedSeason.ChampionTeam())
		assert.Equal(t, SeasonPhaseCompleted, updatedSeason.CurrentPhase())
		assert.Equal(t, SeasonStatusFinished, updatedSeason.CurrentStatus())
	})

	t.Run("Record playoff semifinal should advance winner into final", func(t *testing.T) {
		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusInProgress, SeasonPhasePlayoffs, 0, nil, &PlayoffRulesSnapshot{
			QualificationType: "top_n",
			QualifierCount:    4,
			Rounds: []PlayoffRoundRuleSnapshot{
				{Name: "semifinal", Legs: 1, TiedAggregateResolution: "penalties"},
				{Name: "final", Legs: 1, TiedAggregateResolution: "penalties"},
			},
		}, RehydratePlayoffBracket([]PlayoffBracketRoundSnapshot{
			{
				Name:  "semifinal",
				Order: 1,
				Ties: []PlayoffTieSnapshot{
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
						Matches: []MatchSnapshot{
							RehydrateMatch(MatchState{ID: "semi-leg-1", PlayoffTieID: "semi-1", MatchOrder: 1, HomeTeamID: "team-1", AwayTeamID: "team-4", Status: MatchStatusScheduled}).Snapshot(),
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
						Matches: []MatchSnapshot{
							RehydrateMatch(MatchState{ID: "semi-leg-2", PlayoffTieID: "semi-2", MatchOrder: 1, HomeTeamID: "team-2", AwayTeamID: "team-3", Status: MatchStatusScheduled}).Snapshot(),
						},
					},
				},
			},
			{
				Name:  "final",
				Order: 2,
				Ties: []PlayoffTieSnapshot{
					{
						ID:         "final-1",
						RoundName:  "final",
						RoundOrder: 2,
						SlotOrder:  1,
						Status:     "pending",
						Matches:    []MatchSnapshot{},
					},
				},
			},
		}), nil)

		updatedSeason, err := season.RecordPlayoffMatchScore("semi-1", "semi-leg-1", 2, 0)

		assert.NoError(t, err)
		winnerTie, ok := updatedSeason.FindPlayoffTie("semi-1")
		assert.True(t, ok)
		assert.NotNil(t, winnerTie.WinnerTeamID)
		assert.Equal(t, "team-1", *winnerTie.WinnerTeamID)
		finalRound, ok := updatedSeason.FindPlayoffRound(2)
		assert.True(t, ok)
		assert.Equal(t, "team-1", finalRound.Ties[0].HomeTeamID)
		assert.Empty(t, finalRound.Ties[0].Matches)
	})

	t.Run("Record second playoff semifinal should advance winner into final away slot and create final matches", func(t *testing.T) {
		semiWinner := func() *string { v := "team-1"; return &v }()
		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusInProgress, SeasonPhasePlayoffs, 0, nil, &PlayoffRulesSnapshot{
			QualificationType: "top_n",
			QualifierCount:    4,
			Rounds: []PlayoffRoundRuleSnapshot{
				{Name: "semifinal", Legs: 1, TiedAggregateResolution: "higher_seed_advances"},
				{Name: "final", Legs: 1, TiedAggregateResolution: "clear_winner_required"},
			},
		}, RehydratePlayoffBracket([]PlayoffBracketRoundSnapshot{
			{
				Name:  "semifinal",
				Order: 1,
				Ties: []PlayoffTieSnapshot{
					{
						ID:           "semi-1",
						RoundName:    "semifinal",
						RoundOrder:   1,
						SlotOrder:    1,
						HomeSeed:     1,
						AwaySeed:     4,
						HomeTeamID:   "team-1",
						AwayTeamID:   "team-4",
						Status:       "finished",
						WinnerTeamID: semiWinner,
						Matches: []MatchSnapshot{
							RehydrateMatch(MatchState{ID: "semi-leg-1", PlayoffTieID: "semi-1", MatchOrder: 1, HomeTeamID: "team-1", AwayTeamID: "team-4", Status: MatchStatusFinished, HomeTeamScore: 2, AwayTeamScore: 0}).Snapshot(),
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
						Matches: []MatchSnapshot{
							RehydrateMatch(MatchState{ID: "semi-leg-2", PlayoffTieID: "semi-2", MatchOrder: 1, HomeTeamID: "team-2", AwayTeamID: "team-3", Status: MatchStatusScheduled}).Snapshot(),
						},
					},
				},
			},
			{
				Name:  "final",
				Order: 2,
				Ties: []PlayoffTieSnapshot{
					{
						ID:         "final-1",
						RoundName:  "final",
						RoundOrder: 2,
						SlotOrder:  1,
						HomeTeamID: "team-1",
						HomeSeed:   1,
						Status:     "pending",
						Matches:    []MatchSnapshot{},
					},
				},
			},
		}), nil)

		updatedSeason, err := season.RecordPlayoffMatchScore("semi-2", "semi-leg-2", 1, 0)

		assert.NoError(t, err)
		winnerTie, ok := updatedSeason.FindPlayoffTie("semi-2")
		assert.True(t, ok)
		assert.NotNil(t, winnerTie.WinnerTeamID)
		assert.Equal(t, "team-2", *winnerTie.WinnerTeamID)
		finalRound, ok := updatedSeason.FindPlayoffRound(2)
		assert.True(t, ok)
		assert.Equal(t, "team-2", finalRound.Ties[0].AwayTeamID)
		assert.Equal(t, 2, finalRound.Ties[0].AwaySeed)
		assert.Equal(t, "ready", finalRound.Ties[0].Status)
		assert.Len(t, finalRound.Ties[0].Matches, 1)
	})

	t.Run("Record tied playoff semifinal should advance higher seed", func(t *testing.T) {
		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusInProgress, SeasonPhasePlayoffs, 0, nil, &PlayoffRulesSnapshot{
			QualificationType: "top_n",
			QualifierCount:    4,
			Rounds: []PlayoffRoundRuleSnapshot{
				{Name: "semifinal", Legs: 2, HigherSeedHostsSecondLeg: true, TiedAggregateResolution: "higher_seed_advances"},
				{Name: "final", Legs: 1, TiedAggregateResolution: "clear_winner_required"},
			},
		}, RehydratePlayoffBracket([]PlayoffBracketRoundSnapshot{
			{
				Name:  "semifinal",
				Order: 1,
				Ties: []PlayoffTieSnapshot{
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
						Matches: []MatchSnapshot{
							RehydrateMatch(MatchState{ID: "semi-leg-1", PlayoffTieID: "semi-1", MatchOrder: 1, HomeTeamID: "team-4", AwayTeamID: "team-1", Status: MatchStatusFinished, HomeTeamScore: 2, AwayTeamScore: 1}).Snapshot(),
							RehydrateMatch(MatchState{ID: "semi-leg-2", PlayoffTieID: "semi-1", MatchOrder: 2, HomeTeamID: "team-1", AwayTeamID: "team-4", Status: MatchStatusScheduled}).Snapshot(),
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
				Ties: []PlayoffTieSnapshot{
					{
						ID:         "final-1",
						RoundName:  "final",
						RoundOrder: 2,
						SlotOrder:  1,
						Status:     "pending",
					},
				},
			},
		}), nil)

		updatedSeason, err := season.RecordPlayoffMatchScore("semi-1", "semi-leg-2", 1, 0)

		assert.NoError(t, err)
		tie, ok := updatedSeason.FindPlayoffTie("semi-1")
		assert.True(t, ok)
		assert.NotNil(t, tie.WinnerTeamID)
		assert.Equal(t, "team-1", *tie.WinnerTeamID)
		assert.Equal(t, "finished", tie.Status)
		finalRound, ok := updatedSeason.FindPlayoffRound(2)
		assert.True(t, ok)
		assert.Equal(t, "team-1", finalRound.Ties[0].HomeTeamID)
	})

	t.Run("Record tied playoff final should fail because final must have winner", func(t *testing.T) {
		season := rehydratedSeasonForTest("id-1", "league-id-1", "Test League", SeasonStatusInProgress, SeasonPhasePlayoffs, 0, nil, &PlayoffRulesSnapshot{
			QualificationType: "top_n",
			QualifierCount:    2,
			Rounds: []PlayoffRoundRuleSnapshot{
				{Name: "final", Legs: 1, TiedAggregateResolution: "clear_winner_required"},
			},
		}, RehydratePlayoffBracket([]PlayoffBracketRoundSnapshot{
			{
				Name:  "final",
				Order: 1,
				Ties: []PlayoffTieSnapshot{
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
						Matches: []MatchSnapshot{
							RehydrateMatch(MatchState{ID: "leg-final", PlayoffTieID: "tie-final", MatchOrder: 1, HomeTeamID: "team-1", AwayTeamID: "team-2", Status: MatchStatusScheduled}).Snapshot(),
						},
					},
				},
			},
		}), nil)

		updatedSeason, err := season.RecordPlayoffMatchScore("tie-final", "leg-final", 1, 1)

		assert.Error(t, err)
		assert.Nil(t, updatedSeason)
	})

}

func TestRehydrateSeason_DefensivelyCopiesNestedState(t *testing.T) {
	winnerTeamID := "team-1"
	matchStatus := MatchStatusInProgress
	playoffMatchStatus := MatchStatusScheduled

	rounds := []Round{
		RehydrateRound(1, []Match{
			RehydrateMatch(MatchState{
				ID:         "match-1",
				HomeTeamID: "team-1",
				AwayTeamID: "team-2",
				Status:     matchStatus,
			}),
		}),
	}

	rules := &PlayoffRulesSnapshot{
		QualificationType: "top_n",
		QualifierCount:    2,
		Rounds: []PlayoffRoundRuleSnapshot{
			{Name: "final", Legs: 1, TiedAggregateResolution: "clear_winner_required"},
		},
	}

	bracket := &PlayoffBracketSnapshot{
		Rounds: []PlayoffBracketRoundSnapshot{
			{
				Name:  "final",
				Order: 1,
				Ties: []PlayoffTieSnapshot{
					{
						ID:           "tie-1",
						RoundName:    "final",
						RoundOrder:   1,
						SlotOrder:    1,
						HomeSeed:     1,
						AwaySeed:     2,
						HomeTeamID:   "team-1",
						AwayTeamID:   "team-2",
						Status:       "ready",
						WinnerTeamID: &winnerTeamID,
						Matches: []MatchSnapshot{
							{
								ID:         "playoff-match-1",
								HomeTeamID: "team-1",
								AwayTeamID: "team-2",
								Status:     playoffMatchStatus,
							},
						},
					},
				},
			},
		},
	}

	roundSnapshots := make([]RoundSnapshot, len(rounds))
	for i, round := range rounds {
		roundSnapshots[i] = round.Snapshot()
	}

	season := RehydrateSeasonFromSnapshot(SeasonSnapshot{
		ID:             "season-1",
		LeagueID:       "league-1",
		Name:           "Spring",
		Status:         SeasonStatusInProgress,
		Phase:          SeasonPhasePlayoffs,
		Version:        3,
		Rounds:         roundSnapshots,
		PlayoffRules:   rules,
		PlayoffBracket: bracket,
		ChampionTeamID: &winnerTeamID,
	})

	rounds[0].matches[0].homeTeamID = "mutated-home"
	rules.Rounds[0].Name = "mutated-round"
	bracket.Rounds[0].Ties[0].Matches[0].AwayTeamID = "mutated-away"
	winnerTeamID = "mutated-winner"

	matchSnapshot, ok := season.FindMatchSnapshot("match-1")
	assert.True(t, ok)
	assert.Equal(t, "team-1", matchSnapshot.HomeTeamID)
	playoffRules := season.PlayoffRoundRules()
	assert.Equal(t, "final", playoffRules[0].Name)
	tie, ok := season.FindPlayoffTie("tie-1")
	assert.True(t, ok)
	assert.Equal(t, "team-2", tie.Matches[0].AwayTeamID)
	assert.NotNil(t, season.ChampionTeam())
	assert.Equal(t, "team-1", *season.ChampionTeam())
	assert.NotNil(t, tie.WinnerTeamID)
	assert.Equal(t, "team-1", *tie.WinnerTeamID)
}

func rehydratedSeasonForTest(
	id, leagueID, name string,
	status SeasonStatus,
	phase SeasonPhase,
	version int,
	rounds []Round,
	playoffRules *PlayoffRulesSnapshot,
	playoffBracket *PlayoffBracketSnapshot,
	championTeamID *string,
) Season {
	roundSnapshots := make([]RoundSnapshot, len(rounds))
	for i, round := range rounds {
		roundSnapshots[i] = round.Snapshot()
	}

	return *RehydrateSeasonFromSnapshot(SeasonSnapshot{
		ID:             id,
		LeagueID:       leagueID,
		Name:           name,
		Status:         status,
		Phase:          phase,
		Version:        version,
		Rounds:         roundSnapshots,
		PlayoffRules:   playoffRules,
		PlayoffBracket: playoffBracket,
		ChampionTeamID: championTeamID,
	})
}
