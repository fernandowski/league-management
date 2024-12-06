package domain

import (
	"errors"
	"strings"
)

type Season struct {
	ID             string
	LeagueId       string
	Name           string
	Rounds         []Round
	MatchLocations []MatchLocation
}

func NewSeason(name, leagueID string) (*Season, error) {
	if strings.TrimSpace(leagueID) == "" {
		return nil, errors.New("leagueId cannot be empty")
	}

	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name cannot be empty")
	}

	return &Season{ID: "", Name: name, LeagueId: leagueID}, nil
}

func (s *Season) ScheduleRounds(league League) error {
	if len(league.Memberships) < 2 {
		return errors.New("at least two teams needed to plan season")
	}

	matchUps := generateRoundRobin(league.Memberships)

	s.Rounds = []Round{}

	for roundNumber, pairs := range matchUps {
		round := Round{Matches: []Match{}, RoundNumber: roundNumber + 1}

		for _, pair := range pairs {
			match := Match{
				HomeTeamID: pair[0].TeamID,
				AwayTeamID: pair[1].TeamID,
			}
			round.Matches = append(round.Matches, match)
		}
		s.Rounds = append(s.Rounds, round)
	}
	return nil
}

func generateRoundRobin(leagueMembers []LeagueMembership) [][][]LeagueMembership {
	numTeams := len(leagueMembers)
	members := leagueMembers
	if numTeams%2 != 0 {
		byeMember := LeagueMembership{
			ID:               "bye",
			TeamID:           "bye",
			MemberShipStatus: MembershipActive,
		}
		members = append(members, byeMember)
		numTeams++
	}

	rounds := [][][]LeagueMembership{}

	for i := 0; i < numTeams-1; i++ {
		pairs := [][]LeagueMembership{}
		for j := 0; j < numTeams/2; j++ {
			home := members[j]
			away := members[numTeams-1-j]
			pairs = append(pairs, []LeagueMembership{home, away})
		}

		rounds = append(rounds, pairs)

		members = append([]LeagueMembership{members[0]}, append(members[len(members)-1:], members[1:len(members)-1]...)...)
	}

	return rounds
}
