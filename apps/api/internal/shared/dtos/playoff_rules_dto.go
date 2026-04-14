package dtos

type ConfigurePlayoffRulesDTO struct {
	QualificationType      string                     `json:"qualification_type"`
	QualifierCount         int                        `json:"qualifier_count"`
	ReseedEachRound        bool                       `json:"reseed_each_round"`
	ThirdPlaceMatch        bool                       `json:"third_place_match"`
	AllowAdminSeedOverride bool                       `json:"allow_admin_seed_override"`
	Rounds                 []ConfigurePlayoffRoundDTO `json:"rounds"`
}

type ConfigurePlayoffRoundDTO struct {
	Name                     string `json:"name"`
	Legs                     int    `json:"legs"`
	HigherSeedHostsSecondLeg bool   `json:"higher_seed_hosts_second_leg"`
	TiedAggregateResolution  string `json:"tied_aggregate_resolution"`
}
