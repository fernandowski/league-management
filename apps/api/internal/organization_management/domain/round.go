package domain

type Round struct {
	RoundNumber int
	Matches     []Match
}

func NewRound(roundNumber int) Round {
	return Round{
		RoundNumber: roundNumber,
		Matches:     []Match{},
	}
}

func (r *Round) AddMatches(matches []Match) *Round {
	return &Round{RoundNumber: r.RoundNumber, Matches: matches}
}

func (r *Round) AddMatch(match Match) {
	matches := r.Matches
	r.Matches = append(matches, match)
}
