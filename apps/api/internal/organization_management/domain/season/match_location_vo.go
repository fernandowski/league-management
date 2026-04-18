package season

type MatchLocation struct {
	Name string
}

func NewMatchLocation(name string) MatchLocation {
	return MatchLocation{Name: name}
}
