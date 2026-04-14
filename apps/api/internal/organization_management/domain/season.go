package domain

import (
	"errors"
	"github.com/google/uuid"
	"strings"
)

type SeasonStatus string

const (
	SeasonStatusPending    SeasonStatus = "pending"
	SeasonStatusPlanned    SeasonStatus = "planned"
	SeasonStatusInProgress SeasonStatus = "in_progress"
	SeasonStatusPaused     SeasonStatus = "paused"
	SeasonStatusFinished   SeasonStatus = "finished"
	SeasonStatusUndefined  SeasonStatus = "undefined"
)

type Season struct {
	ID             string
	LeagueId       string
	Name           string
	Status         SeasonStatus
	Version        int
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

	return &Season{ID: uuid.New().String(), Name: name, LeagueId: leagueID, Status: SeasonStatusPending, Version: 0}, nil
}

func (s *Season) ScheduleRounds(league League) error {
	if len(league.Memberships) < 2 {
		return errors.New("at least two teams needed to plan season")
	}
	if s.Status != SeasonStatusPending {
		return errors.New("only season in pending status can be planned")
	}

	matchUps := generateRoundRobin(league.Memberships)

	s.Rounds = []Round{}

	for roundNumber, pairs := range matchUps {
		round := Round{Matches: []Match{}, RoundNumber: roundNumber + 1}

		for _, pair := range pairs {
			match, _ := NewMatch(nil, pair[0].TeamID, pair[1].TeamID)
			round.AddMatch(match)
		}
		s.Rounds = append(s.Rounds, round)
	}

	s.Status = SeasonStatusPlanned
	return nil
}

func (s *Season) Start() (*Season, error) {
	if s.Status != SeasonStatusPlanned {
		return nil, errors.New("season not in correct status must be pending")
	}

	if len(s.Rounds) == 0 {
		return nil, errors.New("cannot start season without rounds planned")
	}

	newSeason := s.copy()
	newSeason.Status = SeasonStatusInProgress
	for matchIndex := range newSeason.Rounds[0].Matches {
		newSeason.Rounds[0].Matches[matchIndex].Status = MatchStatusInProgress
	}

	return newSeason, nil
}

func (s *Season) ChangeMatchScore(matchID string, homeScore, awayScore int) (*Season, error) {
	if s.Status != SeasonStatusInProgress {
		return nil, errors.New("season not in correct status in_progress")
	}

	newSeason := s.copy()

	match, roundIndex, matchIndex := newSeason.findMatch(matchID)
	if match == nil {
		return nil, errors.New("match does not exist")
	}
	if roundIndex != newSeason.findCurrentRoundIndex() {
		return nil, errors.New("cannot change score for match not in current round")
	}
	if match.Status == MatchStatusFinished {
		return nil, errors.New("cannot change score for finished match")
	}
	if match.Status == MatchStatusScheduled {
		match.Status = MatchStatusInProgress
	}

	changedMatch, err := match.ChangeScore(homeScore, awayScore)

	if err != nil {
		return nil, err
	}

	newSeason.Rounds[roundIndex].Matches[matchIndex] = *changedMatch
	return newSeason, nil
}

func (s *Season) CompleteCurrentRound() (*Season, error) {
	if s.Status != SeasonStatusInProgress {
		return nil, errors.New("season not in correct status in_progress")
	}

	newSeason := s.copy()
	currentRoundIndex := newSeason.findCurrentRoundIndex()
	if currentRoundIndex == -1 {
		return nil, errors.New("no current round in progress")
	}

	for matchIndex := range newSeason.Rounds[currentRoundIndex].Matches {
		newSeason.Rounds[currentRoundIndex].Matches[matchIndex].Status = MatchStatusFinished
	}

	nextRoundIndex := currentRoundIndex + 1
	if nextRoundIndex >= len(newSeason.Rounds) {
		newSeason.Status = SeasonStatusFinished
		return newSeason, nil
	}

	for matchIndex := range newSeason.Rounds[nextRoundIndex].Matches {
		newSeason.Rounds[nextRoundIndex].Matches[matchIndex].Status = MatchStatusInProgress
	}

	return newSeason, nil
}

func (s *Season) copy() *Season {
	newSeason := *s

	newRounds := make([]Round, len(s.Rounds))
	for i, round := range s.Rounds {
		newMatches := make([]Match, len(round.Matches))
		copy(newMatches, round.Matches)
		newRounds[i] = Round{
			RoundNumber: round.RoundNumber,
			Matches:     newMatches,
		}
	}

	newSeason.Rounds = newRounds
	newSeason.Version = s.Version

	return &newSeason
}

func (s *Season) findMatch(matchId string) (*Match, int, int) {
	for roundIndex, round := range s.Rounds {
		for matchIndex, match := range round.Matches {
			if match.ID == matchId {
				return &match, roundIndex, matchIndex
			}
		}
	}
	return nil, 0, 0
}

func (s *Season) findCurrentRoundIndex() int {
	for roundIndex, round := range s.Rounds {
		for _, match := range round.Matches {
			if match.Status == MatchStatusInProgress {
				return roundIndex
			}
		}
	}

	for roundIndex, round := range s.Rounds {
		for _, match := range round.Matches {
			if match.Status != MatchStatusFinished {
				return roundIndex
			}
		}
	}

	return -1
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
