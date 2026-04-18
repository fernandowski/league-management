package season

import (
	"errors"
	"github.com/google/uuid"
	"log"
)

type Match struct {
	id               string
	homeTeamID       string
	awayTeamID       string
	homeTeamScore    int
	awayTeamScore    int
	status           MatchStatus
	assignedLocation MatchLocation
	refereeID        string
	playoffTieID     string
	matchOrder       int
}

type MatchState struct {
	ID               string
	HomeTeamID       string
	AwayTeamID       string
	HomeTeamScore    int
	AwayTeamScore    int
	Status           MatchStatus
	AssignedLocation MatchLocation
	RefereeID        string
	PlayoffTieID     string
	MatchOrder       int
}

type MatchSnapshot struct {
	ID               string
	HomeTeamID       string
	AwayTeamID       string
	HomeTeamScore    int
	AwayTeamScore    int
	Status           MatchStatus
	AssignedLocation MatchLocation
	RefereeID        string
	PlayoffTieID     string
	MatchOrder       int
}

type MatchStatus string

const (
	MatchStatusScheduled  MatchStatus = "scheduled"
	MatchStatusInProgress MatchStatus = "in_progress"
	MatchStatusFinished   MatchStatus = "finished"
	MatchStatusSuspended  MatchStatus = "suspended"
	MatchStatusUndefined  MatchStatus = "undefined"
)

func NewMatch(matchID *string, homeTeamID, awayTeamID string) (Match, error) {
	if homeTeamID == "" {
		return Match{}, errors.New("valid match must have a valid home team ID")
	}

	if awayTeamID == "" {
		return Match{}, errors.New("valid match must have a valid away team ID")
	}

	id := uuid.New().String()
	if matchID != nil {
		id = *matchID
	}

	return Match{
		id:            id,
		homeTeamID:    homeTeamID,
		awayTeamID:    awayTeamID,
		homeTeamScore: 0,
		awayTeamScore: 0,
		status:        MatchStatusScheduled,
	}, nil
}

func RehydrateMatch(state MatchState) Match {
	return Match{
		id:               state.ID,
		homeTeamID:       state.HomeTeamID,
		awayTeamID:       state.AwayTeamID,
		homeTeamScore:    state.HomeTeamScore,
		awayTeamScore:    state.AwayTeamScore,
		status:           state.Status,
		assignedLocation: state.AssignedLocation,
		refereeID:        state.RefereeID,
		playoffTieID:     state.PlayoffTieID,
		matchOrder:       state.MatchOrder,
	}
}

func (m Match) Snapshot() MatchSnapshot {
	return MatchSnapshot{
		ID:               m.id,
		HomeTeamID:       m.homeTeamID,
		AwayTeamID:       m.awayTeamID,
		HomeTeamScore:    m.homeTeamScore,
		AwayTeamScore:    m.awayTeamScore,
		Status:           m.status,
		AssignedLocation: m.assignedLocation,
		RefereeID:        m.refereeID,
		PlayoffTieID:     m.playoffTieID,
		MatchOrder:       m.matchOrder,
	}
}

func (m *Match) GetHomeTeam() interface{} {
	if m.homeTeamID == "bye" {
		return nil
	}
	return m.homeTeamID
}

func (m *Match) GetAwayTeam() interface{} {
	if m.awayTeamID == "bye" {
		return nil
	}
	return m.awayTeamID
}

func (m *Match) ChangeScore(homeTeamScore, awayTeamScore int) (*Match, error) {
	if homeTeamScore < 0 {
		return nil, errors.New("score must be greater than 0")
	}

	if awayTeamScore < 0 {
		return nil, errors.New("score must be greater than 0")
	}

	log.Print(m)
	if m.homeTeamID == "bye" || m.awayTeamID == "bye" {
		return nil, errors.New("cannot change score of match for a bye week")
	}
	if m.status != MatchStatusInProgress {
		return nil, errors.New("cannot change score for match not in current round")
	}

	newMatch := *m
	newMatch.awayTeamScore = awayTeamScore
	newMatch.homeTeamScore = homeTeamScore

	return &newMatch, nil
}
