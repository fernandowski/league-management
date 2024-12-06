package domain

type Match struct {
	HomeTeamID       string
	AwayTeamID       string
	AssignedLocation MatchLocation
}

func NewMatch(homeTeamID, awayTeamID string) Match {
	return Match{HomeTeamID: homeTeamID, AwayTeamID: awayTeamID}
}
