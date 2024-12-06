package domain

type Round struct {
	RoundNumber int
	Matches     []Match
}

func NewRound() Round {
	return Round{
		RoundNumber: 0,
		Matches:     []Match{},
	}
}

func (r *Round) AddMatch(match Match) {
	matches := r.Matches
	r.Matches = append(matches, match)
}
