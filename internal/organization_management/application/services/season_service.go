package services

import (
	"errors"
	"league-management/internal/organization_management/domain"
	"league-management/internal/organization_management/domain/domainservices"
	"league-management/internal/organization_management/infrastructure/repositories"
	"league-management/internal/shared/dtos"
)

type SeasonService struct {
	seasonRepository seasonRepository
	leagueRepository seasonLeagueRepository
	organizationRepo seasonOrganizationRepository
}

func NewSeasonService() *SeasonService {
	return &SeasonService{
		seasonRepository: repositories.NewSeasonRepository(),
		leagueRepository: repositories.NewLeagueRepository(),
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

	league, err := ss.leagueRepository.FindById(season.LeagueId)
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

	err = season.ScheduleRounds(*league)
	if err != nil {
		return err
	}

	err = ss.seasonRepository.Save(season)
	if err != nil {
		return err
	}

	return nil
}

func (ss *SeasonService) StartSeason(orgOwnerID, seasonID string) error {
	season, err := ss.seasonRepository.FindByID(seasonID)
	if err != nil {
		return err
	}

	league, err := ss.leagueRepository.FindById(season.LeagueId)
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

	league, err := ss.leagueRepository.FindById(season.LeagueId)
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

func (ss *SeasonService) SeasonDetails(orgOwnerID, seasonID string) (map[string]interface{}, error) {
	season, err := ss.seasonRepository.FindByID(seasonID)
	if err != nil {
		return nil, err
	}

	league, err := ss.leagueRepository.FindById(season.LeagueId)
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

	result, err := ss.seasonRepository.FetchDetails(season.ID)
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

	league, err := ss.leagueRepository.FindById(season.LeagueId)
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

	result, err := ss.seasonRepository.FetchSeasonStandings(season.ID)
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

	league, err := ss.leagueRepository.FindById(season.LeagueId)
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

	result, err := ss.seasonRepository.FetchSeasonMatchUps(season.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
