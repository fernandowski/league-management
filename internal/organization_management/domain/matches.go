package domain

type Match struct {
	HomeTeamID string
	AwayTeamID string
}

func NewMatch(homeTeamID, awayTeamID string) Match {
	return Match{HomeTeamID: homeTeamID, AwayTeamID: awayTeamID}
}
