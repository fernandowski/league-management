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
		assert.Equal(t, 1, firstRound.RoundNumber)
		assert.Len(t, firstRound.Matches, 2)
	})

	t.Run("Season with odd teams without rules", func(t *testing.T) {
		season := Season{
			ID:       "id-1",
			LeagueId: "league-id-1",
			Name:     "Test League",
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
}
