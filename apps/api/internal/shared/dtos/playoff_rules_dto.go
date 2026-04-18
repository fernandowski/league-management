package dtos

type ConfigurePlayoffRulesDTO struct {
	QualificationType string                     `json:"qualification_type"`
	QualifierCount    int                        `json:"qualifier_count"`
	Rounds            []ConfigurePlayoffRoundDTO `json:"rounds"`
}

type ConfigurePlayoffRoundDTO struct {
	Name                     string `json:"name"`
	Legs                     int    `json:"legs"`
	HigherSeedHostsSecondLeg bool   `json:"higher_seed_hosts_second_leg"`
	TiedAggregateResolution  string `json:"tied_aggregate_resolution"`
}
