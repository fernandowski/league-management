package services

import (
	"errors"
	"league-management/internal/organization_management/domain/domainservices"
	"league-management/internal/organization_management/infrastructure/repositories"
	"league-management/internal/shared/dtos"
)

type SeasonService struct {
}

func NewSeasonService() *SeasonService {
	return &SeasonService{}
}

type SearchSeasonDTO struct {
	LeagueId string
	Term     string
	Limit    int
	Offset   int
}

var seasonRepository = repositories.NewSeasonRepository()

func (ss *SeasonService) AddNewSeason(orgOwnerID, leagueID, seasonName string) error {
	league, err := leagueRepository.FindById(leagueID)
	if err != nil {
		return err
	}

	organization, err := organizationRepo.FindById(league.OrganizationId)
	if err != nil {
		return err
	}

	newSeason, err := domainservices.CreateSeason(organization, league, orgOwnerID, seasonName)
	if err != nil {
		return err
	}

	err = seasonRepository.Save(newSeason)
	if err != nil {
		return err
	}

	return nil
}

func (ss *SeasonService) PlanSchedule(orgOwnerID, seasonID string) error {
	season, err := seasonRepository.FindByID(seasonID)
	if err != nil {
		return err
	}

	league, err := leagueRepository.FindById(season.LeagueId)
	if err != nil {
		return err
	}

	organization, err := organizationRepo.FindById(league.OrganizationId)
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

	err = seasonRepository.Save(season)
	if err != nil {
		return err
	}

	return nil
}

func (ss *SeasonService) Search(orgOwnerID, leagueID string, searchDTO dtos.SearchSeasonDTO) map[string]interface{} {
	var data, total = seasonRepository.Search(orgOwnerID, leagueID, searchDTO)

	result := make(map[string]interface{})
	result["data"] = data
	result["total"] = total

	return result
}
