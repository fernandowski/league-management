package domain

import (
	"errors"
	"github.com/google/uuid"
	"strings"
)

type SeasonStatus string
type SeasonPhase string

type PlayoffRules struct {
	qualificationType string
	qualifierCount    int
	rounds            []PlayoffRoundRule
}

type PlayoffRoundRule struct {
	name                     string
	legs                     int
	higherSeedHostsSecondLeg bool
	tiedAggregateResolution  string
}

type PlayoffQualifiedTeam struct {
	TeamID string
	Seed   int
}

type PlayoffBracket struct {
	rounds []PlayoffBracketRound
}

type PlayoffBracketRound struct {
	name  string
	order int
	ties  []PlayoffTie
}

type PlayoffTie struct {
	id           string
	roundName    string
	roundOrder   int
	slotOrder    int
	homeSeed     int
	awaySeed     int
	homeTeamID   string
	awayTeamID   string
	status       string
	matches      []Match
	winnerTeamID *string
}

type PlayoffRulesSnapshot struct {
	QualificationType string
	QualifierCount    int
	Rounds            []PlayoffRoundRuleSnapshot
}

type PlayoffRoundRuleSnapshot struct {
	Name                     string
	Legs                     int
	HigherSeedHostsSecondLeg bool
	TiedAggregateResolution  string
}

type PlayoffBracketSnapshot struct {
	Rounds []PlayoffBracketRoundSnapshot
}

type PlayoffBracketRoundSnapshot struct {
	Name  string
	Order int
	Ties  []PlayoffTieSnapshot
}

type PlayoffTieSnapshot struct {
	ID           string
	RoundName    string
	RoundOrder   int
	SlotOrder    int
	HomeSeed     int
	AwaySeed     int
	HomeTeamID   string
	AwayTeamID   string
	Status       string
	Matches      []MatchSnapshot
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
	id             string
	leagueID       string
	name           string
	status         SeasonStatus
	phase          SeasonPhase
	version        int
	rounds         []Round
	matchLocations []MatchLocation
	playoffRules   *PlayoffRules
	playoffBracket *PlayoffBracket
	championTeamID *string
}

type SeasonSnapshot struct {
	ID             string
	LeagueID       string
	Name           string
	Status         SeasonStatus
	Phase          SeasonPhase
	Version        int
	Rounds         []RoundSnapshot
	PlayoffRules   *PlayoffRulesSnapshot
	PlayoffBracket *PlayoffBracketSnapshot
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
		id:       uuid.New().String(),
		name:     name,
		leagueID: leagueID,
		status:   SeasonStatusPending,
		phase:    SeasonPhaseRegularSeason,
		version:  0,
	}, nil
}

func (s *Season) ScheduleRounds(leagueID string, memberships []LeagueMembership) error {
	if strings.TrimSpace(leagueID) == "" {
		return errors.New("leagueID cannot be empty")
	}
	if s.leagueID != leagueID {
		return errors.New("season does not belong to the provided league")
	}
	if len(memberships) < 2 {
		return errors.New("at least two teams needed to plan season")
	}
	if s.status != SeasonStatusPending {
		return errors.New("only season in pending status can be planned")
	}

	matchUps := generateRoundRobin(copyLeagueMemberships(memberships))

	s.rounds = []Round{}

	for roundNumber, pairs := range matchUps {
		round := Round{matches: []Match{}, roundNumber: roundNumber + 1}

		for _, pair := range pairs {
			match, _ := NewMatch(nil, pair[0].TeamID, pair[1].TeamID)
			round.AddMatch(match)
		}
		s.rounds = append(s.rounds, round)
	}

	s.status = SeasonStatusPlanned
	return nil
}

func copyLeagueMemberships(memberships []LeagueMembership) []LeagueMembership {
	copied := make([]LeagueMembership, len(memberships))
	copy(copied, memberships)
	return copied
}

func (s *Season) Start() (*Season, error) {
	if s.status != SeasonStatusPlanned {
		return nil, errors.New("season not in correct status must be pending")
	}

	if len(s.rounds) == 0 {
		return nil, errors.New("cannot start season without rounds planned")
	}

	newSeason := s.copy()
	newSeason.status = SeasonStatusInProgress
	for matchIndex := range newSeason.rounds[0].matches {
		newSeason.rounds[0].matches[matchIndex].status = MatchStatusInProgress
	}

	return newSeason, nil
}

func (s *Season) ChangeMatchScore(matchID string, homeScore, awayScore int) (*Season, error) {
	if s.status != SeasonStatusInProgress {
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
	if match.status == MatchStatusFinished {
		return nil, errors.New("cannot change score for finished match")
	}
	if match.status == MatchStatusScheduled {
		match.status = MatchStatusInProgress
	}

	changedMatch, err := match.ChangeScore(homeScore, awayScore)

	if err != nil {
		return nil, err
	}

	newSeason.rounds[roundIndex].matches[matchIndex] = *changedMatch
	return newSeason, nil
}

func (s *Season) ChangeMatchScoreByReferee(matchID, refereeID string, homeScore, awayScore int) (*Season, error) {
	if strings.TrimSpace(refereeID) == "" {
		return nil, errors.New("refereeID cannot be empty")
	}

	updatedSeason, err := s.ChangeMatchScore(matchID, homeScore, awayScore)
	if err != nil {
		return nil, err
	}

	match, roundIndex, matchIndex := updatedSeason.findMatch(matchID)
	if match == nil {
		return nil, errors.New("match does not exist")
	}

	updatedMatch := *match
	updatedMatch.refereeID = refereeID
	updatedSeason.rounds[roundIndex].matches[matchIndex] = updatedMatch

	return updatedSeason, nil
}

func (s *Season) CompleteCurrentRound() (*Season, error) {
	if s.status != SeasonStatusInProgress {
		return nil, errors.New("season not in correct status in_progress")
	}

	newSeason := s.copy()
	currentRoundIndex := newSeason.findCurrentRoundIndex()
	if currentRoundIndex == -1 {
		return nil, errors.New("no current round in progress")
	}

	for matchIndex := range newSeason.rounds[currentRoundIndex].matches {
		newSeason.rounds[currentRoundIndex].matches[matchIndex].status = MatchStatusFinished
	}

	nextRoundIndex := currentRoundIndex + 1
	if nextRoundIndex >= len(newSeason.rounds) {
		newSeason.status = SeasonStatusFinished
		return newSeason, nil
	}

	for matchIndex := range newSeason.rounds[nextRoundIndex].matches {
		newSeason.rounds[nextRoundIndex].matches[matchIndex].status = MatchStatusInProgress
	}

	return newSeason, nil
}

func (s *Season) copy() *Season {
	newSeason := *s

	newRounds := make([]Round, len(s.rounds))
	for i, round := range s.rounds {
		newMatches := make([]Match, len(round.matches))
		copy(newMatches, round.matches)
		newRounds[i] = Round{
			roundNumber: round.roundNumber,
			matches:     newMatches,
		}
	}

	newSeason.rounds = newRounds
	newSeason.version = s.version
	newSeason.playoffRules = rehydratePlayoffRulesPtr(s.copyPlayoffRules())
	newSeason.playoffBracket = rehydratePlayoffBracketPtr(s.copyPlayoffBracket())

	return &newSeason
}

func (s *Season) ConfigurePlayoffRules(rules PlayoffRulesSnapshot) (*Season, error) {
	if s.phase == SeasonPhaseCompleted {
		return nil, errors.New("cannot configure playoff rules after playoffs have started")
	}
	if s.playoffsHaveStarted() {
		return nil, errors.New("cannot configure playoff rules after the playoffs have started")
	}

	if err := rules.Validate(); err != nil {
		return nil, err
	}

	newSeason := s.copy()
	rehydratedRules := RehydratePlayoffRules(rules)
	newSeason.playoffRules = &rehydratedRules
	if newSeason.playoffBracket != nil {
		newSeason.playoffBracket = nil
		newSeason.phase = SeasonPhaseRegularSeason
		newSeason.status = SeasonStatusFinished
	}
	return newSeason, nil
}

func (s *Season) playoffsHaveStarted() bool {
	if s.playoffBracket == nil {
		return false
	}

	for _, round := range s.playoffBracket.rounds {
		for _, tie := range round.ties {
			for _, match := range tie.matches {
				if match.status != MatchStatusScheduled {
					return true
				}
			}
		}
	}

	return false
}

func (s *Season) playoffBracketIsReady() bool {
	if s.playoffBracket == nil || len(s.playoffBracket.rounds) == 0 {
		return false
	}

	firstRound := s.playoffBracket.rounds[0]
	if len(firstRound.ties) == 0 {
		return false
	}

	for _, tie := range firstRound.ties {
		if tie.homeTeamID == "" || tie.awayTeamID == "" || len(tie.matches) == 0 {
			return false
		}
	}

	return true
}

func (s *Season) HasStartedPlayoffMatches() bool {
	return s.playoffsHaveStarted()
}

func (s *Season) HasUsablePlayoffBracket() bool {
	return s.playoffBracketIsReady()
}

func (s *Season) HasPlayoffRules() bool {
	return s.playoffRules != nil
}

func (s *Season) HasPlayoffBracket() bool {
	return s.playoffBracket != nil
}

func (s *Season) PlayoffQualificationType() (string, bool) {
	if s.playoffRules == nil {
		return "", false
	}

	return s.playoffRules.qualificationType, true
}

func (s *Season) PlayoffQualifierCount() (int, bool) {
	if s.playoffRules == nil {
		return 0, false
	}

	return s.playoffRules.qualifierCount, true
}

func (s *Season) PlayoffRoundRules() []PlayoffRoundRuleSnapshot {
	if s.playoffRules == nil {
		return nil
	}

	rounds := make([]PlayoffRoundRuleSnapshot, len(s.playoffRules.rounds))
	for i, round := range s.playoffRules.rounds {
		rounds[i] = round.Snapshot()
	}

	return rounds
}

func (s *Season) PlayoffBracketRounds() []PlayoffBracketRoundSnapshot {
	if s.playoffBracket == nil {
		return nil
	}

	rounds := make([]PlayoffBracketRoundSnapshot, len(s.playoffBracket.rounds))
	for i, round := range s.playoffBracket.rounds {
		rounds[i] = round.Snapshot()
	}

	return rounds
}

func (s *Season) FindPlayoffRound(roundOrder int) (*PlayoffBracketRoundSnapshot, bool) {
	if s.playoffBracket == nil {
		return nil, false
	}

	for _, round := range s.playoffBracket.rounds {
		if round.order == roundOrder {
			snapshot := round.Snapshot()
			return &snapshot, true
		}
	}

	return nil, false
}

func (s *Season) FindPlayoffTie(tieID string) (*PlayoffTieSnapshot, bool) {
	if s.playoffBracket == nil {
		return nil, false
	}

	for _, round := range s.playoffBracket.rounds {
		for _, tie := range round.ties {
			if tie.id == tieID {
				snapshot := tie.Snapshot()
				return &snapshot, true
			}
		}
	}

	return nil, false
}

func (s *Season) SeasonID() string {
	return s.id
}

func (s *Season) LeagueID() string {
	return s.leagueID
}

func (s *Season) SeasonName() string {
	return s.name
}

func (s *Season) CurrentStatus() SeasonStatus {
	return s.status
}

func (s *Season) CurrentPhase() SeasonPhase {
	return s.phase
}

func (s *Season) CurrentVersion() int {
	return s.version
}

func (s *Season) ChampionTeam() *string {
	return copyStringPtr(s.championTeamID)
}

func (s *Season) FindRound(roundNumber int) (*RoundSnapshot, bool) {
	for _, round := range s.rounds {
		if round.roundNumber == roundNumber {
			snapshot := round.Snapshot()
			return &snapshot, true
		}
	}

	return nil, false
}

func (s *Season) CurrentRound() (*RoundSnapshot, bool) {
	currentRoundIndex := s.findCurrentRoundIndex()
	if currentRoundIndex == -1 {
		return nil, false
	}

	snapshot := s.rounds[currentRoundIndex].Snapshot()
	return &snapshot, true
}

func (s *Season) FindMatchSnapshot(matchID string) (*MatchSnapshot, bool) {
	match, _, _ := s.findMatch(matchID)
	if match == nil {
		return nil, false
	}

	snapshot := match.Snapshot()
	return &snapshot, true
}

func (s *Season) RoundCount() int {
	return len(s.rounds)
}

func (s *Season) PlannedRounds() []RoundSnapshot {
	return copyRoundSnapshots(s.rounds)
}

func (s *Season) Rules() *PlayoffRulesSnapshot {
	return s.copyPlayoffRules()
}

func (s *Season) Bracket() *PlayoffBracketSnapshot {
	return s.copyPlayoffBracket()
}

func (s *Season) Snapshot() SeasonSnapshot {
	return SeasonSnapshot{
		ID:             s.id,
		LeagueID:       s.leagueID,
		Name:           s.name,
		Status:         s.status,
		Phase:          s.phase,
		Version:        s.version,
		Rounds:         copyRoundSnapshots(s.rounds),
		PlayoffRules:   s.copyPlayoffRules(),
		PlayoffBracket: s.copyPlayoffBracket(),
		ChampionTeamID: copyStringPtr(s.championTeamID),
	}
}

func (s *Season) ApplyPersistedVersion(version int) {
	s.version = version
}

func RehydratePlayoffBracket(rounds []PlayoffBracketRoundSnapshot) *PlayoffBracketSnapshot {
	return copyPlayoffBracketValue(rehydratePlayoffBracketPtr(&PlayoffBracketSnapshot{Rounds: rounds}))
}

func RehydratePlayoffRules(snapshot PlayoffRulesSnapshot) PlayoffRules {
	rounds := make([]PlayoffRoundRule, len(snapshot.Rounds))
	for i, round := range snapshot.Rounds {
		rounds[i] = PlayoffRoundRule{
			name:                     round.Name,
			legs:                     round.Legs,
			higherSeedHostsSecondLeg: round.HigherSeedHostsSecondLeg,
			tiedAggregateResolution:  round.TiedAggregateResolution,
		}
	}

	return PlayoffRules{
		qualificationType: snapshot.QualificationType,
		qualifierCount:    snapshot.QualifierCount,
		rounds:            rounds,
	}
}

func RehydrateSeasonFromSnapshot(snapshot SeasonSnapshot) *Season {
	rounds := make([]Round, len(snapshot.Rounds))
	for i, round := range snapshot.Rounds {
		matches := make([]Match, len(round.Matches))
		for j, match := range round.Matches {
			matches[j] = RehydrateMatch(MatchState{
				ID:               match.ID,
				HomeTeamID:       match.HomeTeamID,
				AwayTeamID:       match.AwayTeamID,
				HomeTeamScore:    match.HomeTeamScore,
				AwayTeamScore:    match.AwayTeamScore,
				Status:           match.Status,
				AssignedLocation: match.AssignedLocation,
				RefereeID:        match.RefereeID,
				PlayoffTieID:     match.PlayoffTieID,
				MatchOrder:       match.MatchOrder,
			})
		}
		rounds[i] = RehydrateRound(round.RoundNumber, matches)
	}

	return &Season{
		id:             snapshot.ID,
		leagueID:       snapshot.LeagueID,
		name:           snapshot.Name,
		status:         snapshot.Status,
		phase:          snapshot.Phase,
		version:        snapshot.Version,
		rounds:         rounds,
		playoffRules:   rehydratePlayoffRulesPtr(snapshot.PlayoffRules),
		playoffBracket: rehydratePlayoffBracketPtr(snapshot.PlayoffBracket),
		championTeamID: copyStringPtr(snapshot.ChampionTeamID),
	}
}

func rehydratePlayoffRulesPtr(snapshot *PlayoffRulesSnapshot) *PlayoffRules {
	if snapshot == nil {
		return nil
	}

	rules := RehydratePlayoffRules(*snapshot)
	return &rules
}

func rehydratePlayoffBracketPtr(snapshot *PlayoffBracketSnapshot) *PlayoffBracket {
	if snapshot == nil {
		return nil
	}

	rounds := make([]PlayoffBracketRound, len(snapshot.Rounds))
	for i, round := range snapshot.Rounds {
		ties := make([]PlayoffTie, len(round.Ties))
		for j, tie := range round.Ties {
			matches := make([]Match, len(tie.Matches))
			for k, match := range tie.Matches {
				matches[k] = RehydrateMatch(MatchState{
					ID:               match.ID,
					HomeTeamID:       match.HomeTeamID,
					AwayTeamID:       match.AwayTeamID,
					HomeTeamScore:    match.HomeTeamScore,
					AwayTeamScore:    match.AwayTeamScore,
					Status:           match.Status,
					AssignedLocation: match.AssignedLocation,
					RefereeID:        match.RefereeID,
					PlayoffTieID:     match.PlayoffTieID,
					MatchOrder:       match.MatchOrder,
				})
			}

			ties[j] = PlayoffTie{
				id:           tie.ID,
				roundName:    tie.RoundName,
				roundOrder:   tie.RoundOrder,
				slotOrder:    tie.SlotOrder,
				homeSeed:     tie.HomeSeed,
				awaySeed:     tie.AwaySeed,
				homeTeamID:   tie.HomeTeamID,
				awayTeamID:   tie.AwayTeamID,
				status:       tie.Status,
				matches:      matches,
				winnerTeamID: copyStringPtr(tie.WinnerTeamID),
			}
		}

		rounds[i] = PlayoffBracketRound{
			name:  round.Name,
			order: round.Order,
			ties:  ties,
		}
	}

	return &PlayoffBracket{rounds: rounds}
}

func copyRounds(rounds []Round) []Round {
	copiedRounds := make([]Round, len(rounds))
	for i, round := range rounds {
		copiedMatches := make([]Match, len(round.matches))
		copy(copiedMatches, round.matches)
		copiedRounds[i] = Round{
			roundNumber: round.roundNumber,
			matches:     copiedMatches,
		}
	}

	return copiedRounds
}

func copyRoundSnapshots(rounds []Round) []RoundSnapshot {
	copied := make([]RoundSnapshot, len(rounds))
	for i, round := range rounds {
		copied[i] = round.Snapshot()
	}
	return copied
}

func (r PlayoffRules) Snapshot() PlayoffRulesSnapshot {
	rounds := make([]PlayoffRoundRuleSnapshot, len(r.rounds))
	for i, round := range r.rounds {
		rounds[i] = round.Snapshot()
	}

	return PlayoffRulesSnapshot{
		QualificationType: r.qualificationType,
		QualifierCount:    r.qualifierCount,
		Rounds:            rounds,
	}
}

func (r PlayoffRulesSnapshot) Validate() error {
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

func (r PlayoffRoundRule) Snapshot() PlayoffRoundRuleSnapshot {
	return PlayoffRoundRuleSnapshot{
		Name:                     r.name,
		Legs:                     r.legs,
		HigherSeedHostsSecondLeg: r.higherSeedHostsSecondLeg,
		TiedAggregateResolution:  r.tiedAggregateResolution,
	}
}

func (r PlayoffRoundRuleSnapshot) Validate() error {
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

func (b PlayoffBracket) Snapshot() PlayoffBracketSnapshot {
	rounds := make([]PlayoffBracketRoundSnapshot, len(b.rounds))
	for i, round := range b.rounds {
		rounds[i] = round.Snapshot()
	}

	return PlayoffBracketSnapshot{Rounds: rounds}
}

func (r PlayoffBracketRound) Snapshot() PlayoffBracketRoundSnapshot {
	ties := make([]PlayoffTieSnapshot, len(r.ties))
	for i, tie := range r.ties {
		ties[i] = tie.Snapshot()
	}

	return PlayoffBracketRoundSnapshot{
		Name:  r.name,
		Order: r.order,
		Ties:  ties,
	}
}

func (t PlayoffTie) Snapshot() PlayoffTieSnapshot {
	matches := make([]MatchSnapshot, len(t.matches))
	for i, match := range t.matches {
		matches[i] = match.Snapshot()
	}

	return PlayoffTieSnapshot{
		ID:           t.id,
		RoundName:    t.roundName,
		RoundOrder:   t.roundOrder,
		SlotOrder:    t.slotOrder,
		HomeSeed:     t.homeSeed,
		AwaySeed:     t.awaySeed,
		HomeTeamID:   t.homeTeamID,
		AwayTeamID:   t.awayTeamID,
		Status:       t.status,
		Matches:      matches,
		WinnerTeamID: copyStringPtr(t.winnerTeamID),
	}
}

func copyPlayoffRulesValue(rules *PlayoffRules) *PlayoffRulesSnapshot {
	if rules == nil {
		return nil
	}

	copiedRounds := make([]PlayoffRoundRuleSnapshot, len(rules.rounds))
	for i, round := range rules.rounds {
		copiedRounds[i] = round.Snapshot()
	}

	return &PlayoffRulesSnapshot{
		QualificationType: rules.qualificationType,
		QualifierCount:    rules.qualifierCount,
		Rounds:            copiedRounds,
	}
}

func copyPlayoffBracketValue(bracket *PlayoffBracket) *PlayoffBracketSnapshot {
	if bracket == nil {
		return nil
	}

	copiedRounds := make([]PlayoffBracketRoundSnapshot, len(bracket.rounds))
	for i, round := range bracket.rounds {
		copiedRounds[i] = round.Snapshot()
	}

	return &PlayoffBracketSnapshot{Rounds: copiedRounds}
}

func copyStringPtr(value *string) *string {
	if value == nil {
		return nil
	}

	copied := *value
	return &copied
}

func (s *Season) copyPlayoffRules() *PlayoffRulesSnapshot {
	return copyPlayoffRulesValue(s.playoffRules)
}

func (s *Season) copyPlayoffBracket() *PlayoffBracketSnapshot {
	return copyPlayoffBracketValue(s.playoffBracket)
}

func (s *Season) GeneratePlayoffBracket(qualifiedTeams []PlayoffQualifiedTeam) (*Season, error) {
	if s.playoffRules == nil {
		return nil, errors.New("playoff rules must be configured before bracket generation")
	}
	if s.phase != SeasonPhaseRegularSeason {
		return nil, errors.New("playoff bracket can only be generated before playoffs start")
	}
	if s.status != SeasonStatusFinished {
		return nil, errors.New("regular season must be finished before bracket generation")
	}
	if s.playoffBracket != nil && s.playoffsHaveStarted() {
		return nil, errors.New("playoff bracket already generated")
	}
	if len(qualifiedTeams) != s.playoffRules.qualifierCount {
		return nil, errors.New("qualified teams count does not match playoff rules")
	}

	newSeason := s.copy()
	bracket := &PlayoffBracket{rounds: []PlayoffBracketRound{}}

	for roundIndex, roundRule := range s.playoffRules.rounds {
		bracketRound := PlayoffBracketRound{
			name:  roundRule.name,
			order: roundIndex + 1,
			ties:  []PlayoffTie{},
		}

		tieCount := len(qualifiedTeams) / 2
		if roundIndex > 0 {
			tieCount = len(bracket.rounds[roundIndex-1].ties) / 2
		}
		if tieCount == 0 {
			tieCount = 1
		}

		for slot := 0; slot < tieCount; slot++ {
			tie := PlayoffTie{
				id:         uuid.New().String(),
				roundName:  roundRule.name,
				roundOrder: roundIndex + 1,
				slotOrder:  slot + 1,
				status:     "pending",
			}

			if roundIndex == 0 {
				homeTeam := qualifiedTeams[slot]
				awayTeam := qualifiedTeams[len(qualifiedTeams)-1-slot]
				tie.homeSeed = homeTeam.Seed
				tie.awaySeed = awayTeam.Seed
				tie.homeTeamID = homeTeam.TeamID
				tie.awayTeamID = awayTeam.TeamID
				tie.matches = buildPlayoffMatches(tie, roundRule)
				tie.status = "ready"
			}

			bracketRound.ties = append(bracketRound.ties, tie)
		}

		bracket.rounds = append(bracket.rounds, bracketRound)
	}

	newSeason.playoffBracket = bracket
	newSeason.phase = SeasonPhasePlayoffs
	newSeason.status = SeasonStatusInProgress
	return newSeason, nil
}

func (s *Season) RecordPlayoffMatchScore(tieID, matchID string, homeScore, awayScore int) (*Season, error) {
	if s.phase != SeasonPhasePlayoffs || s.status != SeasonStatusInProgress {
		return nil, errors.New("playoffs are not in progress")
	}
	if s.playoffBracket == nil {
		return nil, errors.New("playoff bracket has not been generated")
	}
	if homeScore < 0 || awayScore < 0 {
		return nil, errors.New("score must be greater than 0")
	}

	newSeason := s.copy()

	for roundIndex := range newSeason.playoffBracket.rounds {
		for tieIndex := range newSeason.playoffBracket.rounds[roundIndex].ties {
			tie := &newSeason.playoffBracket.rounds[roundIndex].ties[tieIndex]
			if tie.id != tieID {
				continue
			}

			for matchIndex := range tie.matches {
				match := &tie.matches[matchIndex]
				if match.id != matchID {
					continue
				}

				if match.status == MatchStatusScheduled {
					match.status = MatchStatusInProgress
				}
				changedMatch, err := match.ChangeScore(homeScore, awayScore)
				if err != nil {
					return nil, err
				}
				changedMatch.status = MatchStatusFinished
				tie.matches[matchIndex] = *changedMatch

				if allPlayoffMatchesFinished(tie.matches) {
					if err := s.resolvePlayoffTie(newSeason, roundIndex, tieIndex); err != nil {
						return nil, err
					}
				} else {
					tie.status = "in_progress"
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
	for matchOrder := 1; matchOrder <= roundRule.legs; matchOrder++ {
		homeTeamID := tie.homeTeamID
		awayTeamID := tie.awayTeamID

		if roundRule.legs == 2 && matchOrder == 2 {
			if roundRule.higherSeedHostsSecondLeg {
				homeTeamID = tie.homeTeamID
				awayTeamID = tie.awayTeamID
			} else {
				homeTeamID = tie.awayTeamID
				awayTeamID = tie.homeTeamID
			}
		} else if matchOrder == 1 && roundRule.legs == 2 {
			homeTeamID = tie.awayTeamID
			awayTeamID = tie.homeTeamID
		}

		matches = append(matches, Match{
			id:            uuid.New().String(),
			playoffTieID:  tie.id,
			matchOrder:    matchOrder,
			homeTeamID:    homeTeamID,
			awayTeamID:    awayTeamID,
			status:        MatchStatusScheduled,
			homeTeamScore: 0,
			awayTeamScore: 0,
		})
	}
	return matches
}

func allPlayoffMatchesFinished(matches []Match) bool {
	for _, match := range matches {
		if match.status != MatchStatusFinished {
			return false
		}
	}
	return true
}

func (s *Season) resolvePlayoffTie(newSeason *Season, roundIndex, tieIndex int) error {
	tie := &newSeason.playoffBracket.rounds[roundIndex].ties[tieIndex]
	homeAggregate := 0
	awayAggregate := 0

	for _, match := range tie.matches {
		if match.homeTeamID == tie.homeTeamID {
			homeAggregate += match.homeTeamScore
			awayAggregate += match.awayTeamScore
		} else {
			homeAggregate += match.awayTeamScore
			awayAggregate += match.homeTeamScore
		}
	}

	if homeAggregate == awayAggregate {
		if roundIndex == len(newSeason.playoffBracket.rounds)-1 {
			return errors.New("final must have a winner")
		}

		winnerTeamID := tie.homeTeamID
		winnerSeed := tie.homeSeed
		if tie.awaySeed < tie.homeSeed {
			winnerTeamID = tie.awayTeamID
			winnerSeed = tie.awaySeed
		}

		tie.winnerTeamID = &winnerTeamID
		tie.status = "finished"
		return s.advancePlayoffTieWinner(newSeason, roundIndex, tie, winnerTeamID, winnerSeed)
	}

	winnerTeamID := tie.homeTeamID
	winnerSeed := tie.homeSeed
	if awayAggregate > homeAggregate {
		winnerTeamID = tie.awayTeamID
		winnerSeed = tie.awaySeed
	}

	tie.winnerTeamID = &winnerTeamID
	tie.status = "finished"
	return s.advancePlayoffTieWinner(newSeason, roundIndex, tie, winnerTeamID, winnerSeed)
}

func (s *Season) advancePlayoffTieWinner(newSeason *Season, roundIndex int, tie *PlayoffTie, winnerTeamID string, winnerSeed int) error {
	if roundIndex == len(newSeason.playoffBracket.rounds)-1 {
		newSeason.championTeamID = &winnerTeamID
		newSeason.status = SeasonStatusFinished
		newSeason.phase = SeasonPhaseCompleted
		return nil
	}

	nextRound := &newSeason.playoffBracket.rounds[roundIndex+1]
	nextTieIndex := (tie.slotOrder - 1) / 2
	if nextTieIndex >= len(nextRound.ties) {
		return nil
	}

	nextTie := &nextRound.ties[nextTieIndex]
	if tie.slotOrder%2 == 1 {
		nextTie.homeTeamID = winnerTeamID
		nextTie.homeSeed = winnerSeed
	} else {
		nextTie.awayTeamID = winnerTeamID
		nextTie.awaySeed = winnerSeed
	}

	if nextTie.homeTeamID != "" && nextTie.awayTeamID != "" && len(nextTie.matches) == 0 {
		roundRule := s.playoffRules.rounds[roundIndex+1]
		nextTie.matches = buildPlayoffMatches(*nextTie, roundRule)
		nextTie.status = "ready"
	}

	return nil
}

func (r PlayoffRules) Validate() error {
	if strings.TrimSpace(r.qualificationType) == "" {
		return errors.New("qualification type is required")
	}
	if r.qualificationType != "top_n" {
		return errors.New("unsupported qualification type")
	}
	if r.qualifierCount < 2 {
		return errors.New("qualifier count must be at least 2")
	}
	if len(r.rounds) == 0 {
		return errors.New("at least one playoff round is required")
	}

	for _, round := range r.rounds {
		if err := round.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (r PlayoffRoundRule) Validate() error {
	if strings.TrimSpace(r.name) == "" {
		return errors.New("round name is required")
	}
	if r.legs < 1 || r.legs > 2 {
		return errors.New("round legs must be 1 or 2")
	}
	if strings.TrimSpace(r.tiedAggregateResolution) == "" {
		return errors.New("tied aggregate resolution is required")
	}

	return nil
}

func (s *Season) findMatch(matchId string) (*Match, int, int) {
	for roundIndex, round := range s.rounds {
		for matchIndex := range round.matches {
			if round.matches[matchIndex].id == matchId {
				return &s.rounds[roundIndex].matches[matchIndex], roundIndex, matchIndex
			}
		}
	}
	return nil, 0, 0
}

func (s *Season) findCurrentRoundIndex() int {
	for roundIndex, round := range s.rounds {
		for _, match := range round.matches {
			if match.status == MatchStatusInProgress {
				return roundIndex
			}
		}
	}

	for roundIndex, round := range s.rounds {
		for _, match := range round.matches {
			if match.status != MatchStatusFinished {
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
