package domain

import (
	"errors"
	"github.com/google/uuid"
	"strings"
)

type SeasonStatus string
type SeasonPhase string

type PlayoffRules struct {
	QualificationType      string
	QualifierCount         int
	ReseedEachRound        bool
	ThirdPlaceMatch        bool
	AllowAdminSeedOverride bool
	Rounds                 []PlayoffRoundRule
}

type PlayoffRoundRule struct {
	Name                     string
	Legs                     int
	HigherSeedHostsSecondLeg bool
	TiedAggregateResolution  string
}

type PlayoffQualifiedTeam struct {
	TeamID string
	Seed   int
}

type PlayoffBracket struct {
	Rounds []PlayoffBracketRound
}

type PlayoffBracketRound struct {
	Name  string
	Order int
	Ties  []PlayoffTie
}

type PlayoffTie struct {
	ID           string
	RoundName    string
	RoundOrder   int
	SlotOrder    int
	HomeSeed     int
	AwaySeed     int
	HomeTeamID   string
	AwayTeamID   string
	Status       string
	Matches      []Match
	WinnerTeamID *string
}

const (
	SeasonStatusPending    SeasonStatus = "pending"
	SeasonStatusPlanned    SeasonStatus = "planned"
	SeasonStatusInProgress SeasonStatus = "in_progress"
	SeasonStatusPaused     SeasonStatus = "paused"
	SeasonStatusFinished   SeasonStatus = "finished"
	SeasonStatusUndefined  SeasonStatus = "undefined"
)

const (
	SeasonPhaseRegularSeason SeasonPhase = "regular_season"
	SeasonPhasePlayoffs      SeasonPhase = "playoffs"
	SeasonPhaseCompleted     SeasonPhase = "completed"
)

type Season struct {
	ID             string
	LeagueId       string
	Name           string
	Status         SeasonStatus
	Phase          SeasonPhase
	Version        int
	Rounds         []Round
	MatchLocations []MatchLocation
	PlayoffRules   *PlayoffRules
	PlayoffBracket *PlayoffBracket
	ChampionTeamID *string
}

func NewSeason(name, leagueID string) (*Season, error) {
	if strings.TrimSpace(leagueID) == "" {
		return nil, errors.New("leagueId cannot be empty")
	}

	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name cannot be empty")
	}

	return &Season{
		ID:       uuid.New().String(),
		Name:     name,
		LeagueId: leagueID,
		Status:   SeasonStatusPending,
		Phase:    SeasonPhaseRegularSeason,
		Version:  0,
	}, nil
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
	newSeason.PlayoffRules = s.copyPlayoffRules()
	newSeason.PlayoffBracket = s.copyPlayoffBracket()

	return &newSeason
}

func (s *Season) ConfigurePlayoffRules(rules PlayoffRules) (*Season, error) {
	if s.Phase == SeasonPhaseCompleted {
		return nil, errors.New("cannot configure playoff rules after playoffs have started")
	}
	if s.playoffsHaveStarted() {
		return nil, errors.New("cannot configure playoff rules after the playoffs have started")
	}

	if err := rules.Validate(); err != nil {
		return nil, err
	}

	newSeason := s.copy()
	newSeason.PlayoffRules = &rules
	if newSeason.PlayoffBracket != nil {
		newSeason.PlayoffBracket = nil
		newSeason.Phase = SeasonPhaseRegularSeason
		newSeason.Status = SeasonStatusFinished
	}
	return newSeason, nil
}

func (s *Season) playoffsHaveStarted() bool {
	if s.PlayoffBracket == nil {
		return false
	}

	for _, round := range s.PlayoffBracket.Rounds {
		for _, tie := range round.Ties {
			for _, match := range tie.Matches {
				if match.Status != MatchStatusScheduled {
					return true
				}
			}
		}
	}

	return false
}

func (s *Season) playoffBracketIsReady() bool {
	if s.PlayoffBracket == nil || len(s.PlayoffBracket.Rounds) == 0 {
		return false
	}

	firstRound := s.PlayoffBracket.Rounds[0]
	if len(firstRound.Ties) == 0 {
		return false
	}

	for _, tie := range firstRound.Ties {
		if tie.HomeTeamID == "" || tie.AwayTeamID == "" || len(tie.Matches) == 0 {
			return false
		}
	}

	return true
}

func (s *Season) copyPlayoffRules() *PlayoffRules {
	if s.PlayoffRules == nil {
		return nil
	}

	copiedRounds := make([]PlayoffRoundRule, len(s.PlayoffRules.Rounds))
	copy(copiedRounds, s.PlayoffRules.Rounds)

	copiedRules := *s.PlayoffRules
	copiedRules.Rounds = copiedRounds

	return &copiedRules
}

func (s *Season) copyPlayoffBracket() *PlayoffBracket {
	if s.PlayoffBracket == nil {
		return nil
	}

	copiedRounds := make([]PlayoffBracketRound, len(s.PlayoffBracket.Rounds))
	for i, round := range s.PlayoffBracket.Rounds {
		copiedTies := make([]PlayoffTie, len(round.Ties))
		for j, tie := range round.Ties {
			copiedMatches := make([]Match, len(tie.Matches))
			copy(copiedMatches, tie.Matches)
			copiedTies[j] = tie
			copiedTies[j].Matches = copiedMatches
		}
		copiedRounds[i] = PlayoffBracketRound{
			Name:  round.Name,
			Order: round.Order,
			Ties:  copiedTies,
		}
	}

	return &PlayoffBracket{Rounds: copiedRounds}
}

func (s *Season) GeneratePlayoffBracket(qualifiedTeams []PlayoffQualifiedTeam) (*Season, error) {
	if s.PlayoffRules == nil {
		return nil, errors.New("playoff rules must be configured before bracket generation")
	}
	if s.Phase != SeasonPhaseRegularSeason {
		return nil, errors.New("playoff bracket can only be generated before playoffs start")
	}
	if s.Status != SeasonStatusFinished {
		return nil, errors.New("regular season must be finished before bracket generation")
	}
	if s.PlayoffBracket != nil && s.playoffsHaveStarted() {
		return nil, errors.New("playoff bracket already generated")
	}
	if len(qualifiedTeams) != s.PlayoffRules.QualifierCount {
		return nil, errors.New("qualified teams count does not match playoff rules")
	}

	newSeason := s.copy()
	bracket := &PlayoffBracket{Rounds: []PlayoffBracketRound{}}

	for roundIndex, roundRule := range s.PlayoffRules.Rounds {
		bracketRound := PlayoffBracketRound{
			Name:  roundRule.Name,
			Order: roundIndex + 1,
			Ties:  []PlayoffTie{},
		}

		tieCount := len(qualifiedTeams) / 2
		if roundIndex > 0 {
			tieCount = len(bracket.Rounds[roundIndex-1].Ties) / 2
		}
		if tieCount == 0 {
			tieCount = 1
		}

		for slot := 0; slot < tieCount; slot++ {
			tie := PlayoffTie{
				ID:         uuid.New().String(),
				RoundName:  roundRule.Name,
				RoundOrder: roundIndex + 1,
				SlotOrder:  slot + 1,
				Status:     "pending",
			}

			if roundIndex == 0 {
				homeTeam := qualifiedTeams[slot]
				awayTeam := qualifiedTeams[len(qualifiedTeams)-1-slot]
				tie.HomeSeed = homeTeam.Seed
				tie.AwaySeed = awayTeam.Seed
				tie.HomeTeamID = homeTeam.TeamID
				tie.AwayTeamID = awayTeam.TeamID
				tie.Matches = buildPlayoffMatches(tie, roundRule)
				tie.Status = "ready"
			}

			bracketRound.Ties = append(bracketRound.Ties, tie)
		}

		bracket.Rounds = append(bracket.Rounds, bracketRound)
	}

	newSeason.PlayoffBracket = bracket
	newSeason.Phase = SeasonPhasePlayoffs
	newSeason.Status = SeasonStatusInProgress
	return newSeason, nil
}

func (s *Season) RecordPlayoffMatchScore(tieID, matchID string, homeScore, awayScore int) (*Season, error) {
	if s.Phase != SeasonPhasePlayoffs || s.Status != SeasonStatusInProgress {
		return nil, errors.New("playoffs are not in progress")
	}
	if s.PlayoffBracket == nil {
		return nil, errors.New("playoff bracket has not been generated")
	}
	if homeScore < 0 || awayScore < 0 {
		return nil, errors.New("score must be greater than 0")
	}

	newSeason := s.copy()

	for roundIndex := range newSeason.PlayoffBracket.Rounds {
		for tieIndex := range newSeason.PlayoffBracket.Rounds[roundIndex].Ties {
			tie := &newSeason.PlayoffBracket.Rounds[roundIndex].Ties[tieIndex]
			if tie.ID != tieID {
				continue
			}

			for matchIndex := range tie.Matches {
				match := &tie.Matches[matchIndex]
				if match.ID != matchID {
					continue
				}

				if match.Status == MatchStatusScheduled {
					match.Status = MatchStatusInProgress
				}
				changedMatch, err := match.ChangeScore(homeScore, awayScore)
				if err != nil {
					return nil, err
				}
				changedMatch.Status = MatchStatusFinished
				tie.Matches[matchIndex] = *changedMatch

				if allPlayoffMatchesFinished(tie.Matches) {
					if err := s.resolvePlayoffTie(newSeason, roundIndex, tieIndex); err != nil {
						return nil, err
					}
				} else {
					tie.Status = "in_progress"
				}

				return newSeason, nil
			}

			return nil, errors.New("playoff match not found")
		}
	}

	return nil, errors.New("playoff tie not found")
}

func buildPlayoffMatches(tie PlayoffTie, roundRule PlayoffRoundRule) []Match {
	matches := []Match{}
	for matchOrder := 1; matchOrder <= roundRule.Legs; matchOrder++ {
		homeTeamID := tie.HomeTeamID
		awayTeamID := tie.AwayTeamID

		if roundRule.Legs == 2 && matchOrder == 2 {
			if roundRule.HigherSeedHostsSecondLeg {
				homeTeamID = tie.HomeTeamID
				awayTeamID = tie.AwayTeamID
			} else {
				homeTeamID = tie.AwayTeamID
				awayTeamID = tie.HomeTeamID
			}
		} else if matchOrder == 1 && roundRule.Legs == 2 {
			homeTeamID = tie.AwayTeamID
			awayTeamID = tie.HomeTeamID
		}

		matches = append(matches, Match{
			ID:            uuid.New().String(),
			PlayoffTieID:  tie.ID,
			MatchOrder:    matchOrder,
			HomeTeamID:    homeTeamID,
			AwayTeamID:    awayTeamID,
			Status:        MatchStatusScheduled,
			HomeTeamScore: 0,
			AwayTeamScore: 0,
		})
	}
	return matches
}

func allPlayoffMatchesFinished(matches []Match) bool {
	for _, match := range matches {
		if match.Status != MatchStatusFinished {
			return false
		}
	}
	return true
}

func (s *Season) resolvePlayoffTie(newSeason *Season, roundIndex, tieIndex int) error {
	tie := &newSeason.PlayoffBracket.Rounds[roundIndex].Ties[tieIndex]
	homeAggregate := 0
	awayAggregate := 0

	for _, match := range tie.Matches {
		if match.HomeTeamID == tie.HomeTeamID {
			homeAggregate += match.HomeTeamScore
			awayAggregate += match.AwayTeamScore
		} else {
			homeAggregate += match.AwayTeamScore
			awayAggregate += match.HomeTeamScore
		}
	}

	if homeAggregate == awayAggregate {
		if roundIndex == len(newSeason.PlayoffBracket.Rounds)-1 {
			return errors.New("final must have a winner")
		}

		winnerTeamID := tie.HomeTeamID
		winnerSeed := tie.HomeSeed
		if tie.AwaySeed < tie.HomeSeed {
			winnerTeamID = tie.AwayTeamID
			winnerSeed = tie.AwaySeed
		}

		tie.WinnerTeamID = &winnerTeamID
		tie.Status = "finished"
		return s.advancePlayoffTieWinner(newSeason, roundIndex, tie, winnerTeamID, winnerSeed)
	}

	winnerTeamID := tie.HomeTeamID
	winnerSeed := tie.HomeSeed
	if awayAggregate > homeAggregate {
		winnerTeamID = tie.AwayTeamID
		winnerSeed = tie.AwaySeed
	}

	tie.WinnerTeamID = &winnerTeamID
	tie.Status = "finished"
	return s.advancePlayoffTieWinner(newSeason, roundIndex, tie, winnerTeamID, winnerSeed)
}

func (s *Season) advancePlayoffTieWinner(newSeason *Season, roundIndex int, tie *PlayoffTie, winnerTeamID string, winnerSeed int) error {
	if roundIndex == len(newSeason.PlayoffBracket.Rounds)-1 {
		newSeason.ChampionTeamID = &winnerTeamID
		newSeason.Status = SeasonStatusFinished
		newSeason.Phase = SeasonPhaseCompleted
		return nil
	}

	nextRound := &newSeason.PlayoffBracket.Rounds[roundIndex+1]
	nextTieIndex := (tie.SlotOrder - 1) / 2
	if nextTieIndex >= len(nextRound.Ties) {
		return nil
	}

	nextTie := &nextRound.Ties[nextTieIndex]
	if tie.SlotOrder%2 == 1 {
		nextTie.HomeTeamID = winnerTeamID
		nextTie.HomeSeed = winnerSeed
	} else {
		nextTie.AwayTeamID = winnerTeamID
		nextTie.AwaySeed = winnerSeed
	}

	if nextTie.HomeTeamID != "" && nextTie.AwayTeamID != "" && len(nextTie.Matches) == 0 {
		roundRule := s.PlayoffRules.Rounds[roundIndex+1]
		nextTie.Matches = buildPlayoffMatches(*nextTie, roundRule)
		nextTie.Status = "ready"
	}

	return nil
}

func (r PlayoffRules) Validate() error {
	if strings.TrimSpace(r.QualificationType) == "" {
		return errors.New("qualification type is required")
	}
	if r.QualificationType != "top_n" {
		return errors.New("unsupported qualification type")
	}
	if r.QualifierCount < 2 {
		return errors.New("qualifier count must be at least 2")
	}
	if len(r.Rounds) == 0 {
		return errors.New("at least one playoff round is required")
	}

	for _, round := range r.Rounds {
		if err := round.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (r PlayoffRoundRule) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("round name is required")
	}
	if r.Legs < 1 || r.Legs > 2 {
		return errors.New("round legs must be 1 or 2")
	}
	if strings.TrimSpace(r.TiedAggregateResolution) == "" {
		return errors.New("tied aggregate resolution is required")
	}

	return nil
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
