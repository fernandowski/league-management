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
