package services

import (
	"errors"
	"league-management/internal/organization_management/domain"
	"league-management/internal/organization_management/domain/domainservices"
	"league-management/internal/shared/dtos"
)

type SeasonService struct {
	seasonRepository seasonRepository
	leagueRepository seasonLeagueRepository
	organizationRepo seasonOrganizationRepository
}

func NewSeasonService(
	seasonRepository seasonRepository,
	leagueRepository seasonLeagueRepository,
	organizationRepo seasonOrganizationRepository,
) *SeasonService {
	return &SeasonService{
		seasonRepository: seasonRepository,
		leagueRepository: leagueRepository,
		organizationRepo: organizationRepo,
	}
}

type SearchSeasonDTO struct {
	LeagueId string
	Term     string
	Limit    int
	Offset   int
}

type seasonRepository interface {
	FindByID(string) (*domain.Season, error)
	Save(*domain.Season) error
	Search(string, string, dtos.SearchSeasonDTO) ([]interface{}, int)
	FetchDetails(string) (map[string]interface{}, error)
	FetchSeasonStandings(string) (map[string]interface{}, error)
	FetchSeasonMatchUps(string) ([]interface{}, error)
	FetchPlayoffBracket(string) (map[string]interface{}, error)
}

type seasonLeagueRepository interface {
	FindById(string) (*domain.League, error)
}

type seasonOrganizationRepository interface {
	FindById(string) (*domain.Organization, error)
}

func (ss *SeasonService) AddNewSeason(orgOwnerID, leagueID, seasonName string) error {
	league, err := ss.leagueRepository.FindById(leagueID)
	if err != nil {
		return err
	}

	organization, err := ss.organizationRepo.FindById(league.OrganizationId)
	if err != nil {
		return err
	}

	newSeason, err := domainservices.CreateSeason(organization, league, orgOwnerID, seasonName)
	if err != nil {
		return err
	}

	err = ss.seasonRepository.Save(newSeason)
	if err != nil {
		return err
	}

	return nil
}

func (ss *SeasonService) PlanSchedule(orgOwnerID, seasonID string) error {
	season, err := ss.seasonRepository.FindByID(seasonID)
	if err != nil {
		return err
	}

	league, err := ss.leagueRepository.FindById(season.LeagueID())
	if err != nil {
		return err
	}

	organization, err := ss.organizationRepo.FindById(league.OrganizationId)
	if err != nil {
		return err
	}

	if !organization.BelongsToOwner(orgOwnerID) {
		return errors.New("only org owner can plan schedule")
	}

	err = season.ScheduleRounds(leagueIDValue(league.Id), activeMemberships(league.Memberships))
	if err != nil {
		return err
	}

	err = ss.seasonRepository.Save(season)
	if err != nil {
		return err
	}

	return nil
}

func activeMemberships(memberships []domain.LeagueMembership) []domain.LeagueMembership {
	active := make([]domain.LeagueMembership, 0, len(memberships))
	for _, membership := range memberships {
		if membership.MemberShipStatus == domain.MembershipActive {
			active = append(active, membership)
		}
	}

	return active
}

func leagueIDValue(id *string) string {
	if id == nil {
		return ""
	}

	return *id
}

func (ss *SeasonService) StartSeason(orgOwnerID, seasonID string) error {
	season, err := ss.seasonRepository.FindByID(seasonID)
	if err != nil {
		return err
	}

	league, err := ss.leagueRepository.FindById(season.LeagueID())
	if err != nil {
		return err
	}

	organization, err := ss.organizationRepo.FindById(league.OrganizationId)
	if err != nil {
		return err
	}

	_, err = domainservices.OrganizationOwnerFromUserId(&orgOwnerID, organization)
	if err != nil {
		return err
	}

	startedSeason, err := season.Start()
	if err != nil {
		return err
	}

	err = ss.seasonRepository.Save(startedSeason)
	if err != nil {
		return err
	}

	return nil
}

func (ss *SeasonService) ChangeMatchUpScore(orgOwnerID, seasonID string, changeScoreDTO dtos.ChangeGameScoreDTO) error {
	season, err := ss.seasonRepository.FindByID(seasonID)
	if err != nil {
		return err
	}

	league, err := ss.leagueRepository.FindById(season.LeagueID())
	if err != nil {
		return err
	}

	organization, err := ss.organizationRepo.FindById(league.OrganizationId)
	if err != nil {
		return err
	}

	_, err = domainservices.OrganizationOwnerFromUserId(&orgOwnerID, organization)
	if err != nil {
		return err
	}

	season, err = season.ChangeMatchScore(changeScoreDTO.MatchID, changeScoreDTO.HomeScore, changeScoreDTO.AwayScore)
	if err != nil {
		return err
	}

	err = ss.seasonRepository.Save(season)
	if err != nil {
		return err
	}

	return nil
}

func (ss *SeasonService) CompleteCurrentRound(orgOwnerID, seasonID string) error {
	season, err := ss.seasonRepository.FindByID(seasonID)
	if err != nil {
		return err
	}

	league, err := ss.leagueRepository.FindById(season.LeagueID())
	if err != nil {
		return err
	}

	organization, err := ss.organizationRepo.FindById(league.OrganizationId)
	if err != nil {
		return err
	}

	_, err = domainservices.OrganizationOwnerFromUserId(&orgOwnerID, organization)
	if err != nil {
		return err
	}

	season, err = season.CompleteCurrentRound()
	if err != nil {
		return err
	}

	err = ss.seasonRepository.Save(season)
	if err != nil {
		return err
	}

	return nil
}

func (ss *SeasonService) SeasonDetails(orgOwnerID, seasonID string) (map[string]interface{}, error) {
	season, err := ss.seasonRepository.FindByID(seasonID)
	if err != nil {
		return nil, err
	}

	league, err := ss.leagueRepository.FindById(season.LeagueID())
	if err != nil {
		return nil, err
	}

	organization, err := ss.organizationRepo.FindById(league.OrganizationId)
	if err != nil {
		return nil, err
	}

	if !organization.BelongsToOwner(orgOwnerID) {
		return nil, errors.New("only org owner can view details")
	}

	result, err := ss.seasonRepository.FetchDetails(season.SeasonID())
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ss *SeasonService) Search(orgOwnerID, leagueID string, searchDTO dtos.SearchSeasonDTO) map[string]interface{} {
	var data, total = ss.seasonRepository.Search(orgOwnerID, leagueID, searchDTO)

	result := make(map[string]interface{})
	result["data"] = data
	result["total"] = total

	return result
}

func (ss *SeasonService) SeasonStandings(orgOwnerID, seasonID string) (map[string]interface{}, error) {
	season, err := ss.seasonRepository.FindByID(seasonID)
	if err != nil {
		return nil, err
	}

	league, err := ss.leagueRepository.FindById(season.LeagueID())
	if err != nil {
		return nil, err
	}

	organization, err := ss.organizationRepo.FindById(league.OrganizationId)
	if err != nil {
		return nil, err
	}

	if !organization.BelongsToOwner(orgOwnerID) {
		return nil, errors.New("only org owner can view details")
	}

	result, err := ss.seasonRepository.FetchSeasonStandings(season.SeasonID())
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (ss *SeasonService) SeasonMatchUps(orgOwnerID, seasonID string) ([]interface{}, error) {
	season, err := ss.seasonRepository.FindByID(seasonID)
	if err != nil {
		return nil, err
	}

	league, err := ss.leagueRepository.FindById(season.LeagueID())
	if err != nil {
		return nil, err
	}

	organization, err := ss.organizationRepo.FindById(league.OrganizationId)
	if err != nil {
		return nil, err
	}

	if !organization.BelongsToOwner(orgOwnerID) {
		return nil, errors.New("only org owner can view details")
	}

	result, err := ss.seasonRepository.FetchSeasonMatchUps(season.SeasonID())
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ss *SeasonService) ConfigurePlayoffRules(orgOwnerID, seasonID string, dto dtos.ConfigurePlayoffRulesDTO) error {
	season, err := ss.seasonRepository.FindByID(seasonID)
	if err != nil {
		return err
	}

	league, err := ss.leagueRepository.FindById(season.LeagueID())
	if err != nil {
		return err
	}

	organization, err := ss.organizationRepo.FindById(league.OrganizationId)
	if err != nil {
		return err
	}

	if !organization.BelongsToOwner(orgOwnerID) {
		return errors.New("only org owner can configure playoff rules")
	}

	rules := domain.PlayoffRulesSnapshot{
		QualificationType: dto.QualificationType,
		QualifierCount:    dto.QualifierCount,
		Rounds:            make([]domain.PlayoffRoundRuleSnapshot, 0, len(dto.Rounds)),
	}

	for _, round := range dto.Rounds {
		rules.Rounds = append(rules.Rounds, domain.PlayoffRoundRuleSnapshot{
			Name:                     round.Name,
			Legs:                     round.Legs,
			HigherSeedHostsSecondLeg: round.HigherSeedHostsSecondLeg,
			TiedAggregateResolution:  round.TiedAggregateResolution,
		})
	}

	updatedSeason, err := season.ConfigurePlayoffRules(rules)
	if err != nil {
		return err
	}

	return ss.seasonRepository.Save(updatedSeason)
}

func (ss *SeasonService) PlayoffRules(orgOwnerID, seasonID string) (map[string]interface{}, error) {
	season, err := ss.seasonRepository.FindByID(seasonID)
	if err != nil {
		return nil, err
	}

	league, err := ss.leagueRepository.FindById(season.LeagueID())
	if err != nil {
		return nil, err
	}

	organization, err := ss.organizationRepo.FindById(league.OrganizationId)
	if err != nil {
		return nil, err
	}

	if !organization.BelongsToOwner(orgOwnerID) {
		return nil, errors.New("only org owner can view playoff rules")
	}

	qualificationType, configured := season.PlayoffQualificationType()
	qualifierCount, _ := season.PlayoffQualifierCount()
	roundRules := season.PlayoffRoundRules()

	result := map[string]interface{}{
		"season_id":         season.SeasonID(),
		"season_status":     season.CurrentStatus(),
		"season_phase":      season.CurrentPhase(),
		"configured":        configured,
		"bracket_generated": season.HasUsablePlayoffBracket(),
		"playoffs_started":  season.HasPlayoffBracket() && season.CurrentPhase() == domain.SeasonPhasePlayoffs && season.CurrentStatus() == domain.SeasonStatusInProgress && season.HasStartedPlayoffMatches(),
		"rules_locked":      season.CurrentPhase() == domain.SeasonPhaseCompleted || season.HasStartedPlayoffMatches(),
		"rules":             nil,
	}

	if configured {
		rounds := make([]map[string]interface{}, 0, len(roundRules))
		for _, round := range roundRules {
			rounds = append(rounds, map[string]interface{}{
				"name":                         round.Name,
				"legs":                         round.Legs,
				"higher_seed_hosts_second_leg": round.HigherSeedHostsSecondLeg,
				"tied_aggregate_resolution":    round.TiedAggregateResolution,
			})
		}

		result["rules"] = map[string]interface{}{
			"qualification_type": qualificationType,
			"qualifier_count":    qualifierCount,
			"rounds":             rounds,
		}
	}

	return result, nil
}

func (ss *SeasonService) PlayoffQualificationPreview(orgOwnerID, seasonID string) (map[string]interface{}, error) {
	season, err := ss.seasonRepository.FindByID(seasonID)
	if err != nil {
		return nil, err
	}

	league, err := ss.leagueRepository.FindById(season.LeagueID())
	if err != nil {
		return nil, err
	}

	organization, err := ss.organizationRepo.FindById(league.OrganizationId)
	if err != nil {
		return nil, err
	}

	if !organization.BelongsToOwner(orgOwnerID) {
		return nil, errors.New("only org owner can view playoff qualification")
	}
	qualifierCount, configured := season.PlayoffQualifierCount()
	if !configured {
		return nil, errors.New("playoff rules must be configured first")
	}

	standings, err := ss.seasonRepository.FetchSeasonStandings(season.SeasonID())
	if err != nil {
		return nil, err
	}

	standingRows, ok := standings["standings"].([]interface{})
	if !ok {
		return nil, errors.New("invalid standings response")
	}

	limit := qualifierCount
	if len(standingRows) < limit {
		limit = len(standingRows)
	}

	qualifiedTeams := make([]interface{}, 0, limit)
	for index := 0; index < limit; index++ {
		row, ok := standingRows[index].(map[string]interface{})
		if !ok {
			continue
		}
		qualifiedTeams = append(qualifiedTeams, map[string]interface{}{
			"team_id":     row["team_id"],
			"team_name":   row["team_name"],
			"seed":        index + 1,
			"rank":        index + 1,
			"points":      row["total_points"],
			"wins":        row["total_wins"],
			"losses":      row["total_losses"],
			"ties":        row["total_ties"],
			"total_goals": row["total_goals"],
		})
	}

	return map[string]interface{}{
		"season_id":       season.SeasonID(),
		"qualifier_count": qualifierCount,
		"qualified_teams": qualifiedTeams,
		"bracket_exists":  season.HasUsablePlayoffBracket(),
	}, nil
}

func (ss *SeasonService) GeneratePlayoffBracket(orgOwnerID, seasonID string) error {
	season, err := ss.seasonRepository.FindByID(seasonID)
	if err != nil {
		return err
	}

	league, err := ss.leagueRepository.FindById(season.LeagueID())
	if err != nil {
		return err
	}

	organization, err := ss.organizationRepo.FindById(league.OrganizationId)
	if err != nil {
		return err
	}

	if !organization.BelongsToOwner(orgOwnerID) {
		return errors.New("only org owner can generate playoff bracket")
	}
	qualifierCount, configured := season.PlayoffQualifierCount()
	if !configured {
		return errors.New("playoff rules must be configured first")
	}

	standings, err := ss.seasonRepository.FetchSeasonStandings(season.SeasonID())
	if err != nil {
		return err
	}

	standingRows, ok := standings["standings"].([]interface{})
	if !ok {
		return errors.New("invalid standings response")
	}
	if len(standingRows) < qualifierCount {
		return errors.New("not enough ranked teams to generate playoff bracket")
	}

	qualifiedTeams := make([]domain.PlayoffQualifiedTeam, 0, qualifierCount)
	for index := 0; index < qualifierCount; index++ {
		row, ok := standingRows[index].(map[string]interface{})
		if !ok {
			return errors.New("invalid standings row")
		}

		teamID, ok := row["team_id"].(string)
		if !ok {
			return errors.New("invalid team_id in standings row")
		}

		qualifiedTeams = append(qualifiedTeams, domain.PlayoffQualifiedTeam{
			TeamID: teamID,
			Seed:   index + 1,
		})
	}

	updatedSeason, err := season.GeneratePlayoffBracket(qualifiedTeams)
	if err != nil {
		return err
	}

	return ss.seasonRepository.Save(updatedSeason)
}

func (ss *SeasonService) PlayoffBracket(orgOwnerID, seasonID string) (map[string]interface{}, error) {
	season, err := ss.seasonRepository.FindByID(seasonID)
	if err != nil {
		return nil, err
	}

	league, err := ss.leagueRepository.FindById(season.LeagueID())
	if err != nil {
		return nil, err
	}

	organization, err := ss.organizationRepo.FindById(league.OrganizationId)
	if err != nil {
		return nil, err
	}

	if !organization.BelongsToOwner(orgOwnerID) {
		return nil, errors.New("only org owner can view playoff bracket")
	}

	return ss.seasonRepository.FetchPlayoffBracket(season.SeasonID())
}

func (ss *SeasonService) RecordPlayoffMatchScore(orgOwnerID, seasonID, tieID, matchID string, dto dtos.ChangeGameScoreDTO) error {
	season, err := ss.seasonRepository.FindByID(seasonID)
	if err != nil {
		return err
	}

	league, err := ss.leagueRepository.FindById(season.LeagueID())
	if err != nil {
		return err
	}

	organization, err := ss.organizationRepo.FindById(league.OrganizationId)
	if err != nil {
		return err
	}

	if !organization.BelongsToOwner(orgOwnerID) {
		return errors.New("only org owner can update playoff scores")
	}

	updatedSeason, err := season.RecordPlayoffMatchScore(tieID, matchID, dto.HomeScore, dto.AwayScore)
	if err != nil {
		return err
	}

	return ss.seasonRepository.Save(updatedSeason)
}
