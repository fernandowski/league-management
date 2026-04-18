package season

type Round struct {
	roundNumber int
	matches     []Match
}

type RoundSnapshot struct {
	RoundNumber int
	Matches     []MatchSnapshot
}

func NewRound(roundNumber int) Round {
	return Round{
		roundNumber: roundNumber,
		matches:     []Match{},
	}
}

func RehydrateRound(roundNumber int, matches []Match) Round {
	copiedMatches := make([]Match, len(matches))
	copy(copiedMatches, matches)

	return Round{
		roundNumber: roundNumber,
		matches:     copiedMatches,
	}
}

func (r Round) Snapshot() RoundSnapshot {
	matches := make([]MatchSnapshot, len(r.matches))
	for i, match := range r.matches {
		matches[i] = match.Snapshot()
	}

	return RoundSnapshot{
		RoundNumber: r.roundNumber,
		Matches:     matches,
	}
}

func (r *Round) AddMatches(matches []Match) *Round {
	round := RehydrateRound(r.roundNumber, matches)
	return &round
}

func (r *Round) AddMatch(match Match) {
	matches := r.matches
	r.matches = append(matches, match)
}
