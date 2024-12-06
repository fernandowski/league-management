package services

import (
	"errors"
	"league-management/internal/organization_management/domain"
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

	if !organization.BelongsToOwner(orgOwnerID) {
		return errors.New("only organization owner allowed to Add a new season")
	}

	newSeason, err := domain.NewSeason(seasonName, *league.Id)
	if err != nil {
		return err
	}

	err = seasonRepository.Save(newSeason)
	if err != nil {
		return err
	}

	return nil
}

func (ss *SeasonService) PlanSchedule() error {
	return nil
}

func (ss *SeasonService) Search(orgOwnerID, leagueID string, searchDTO dtos.SearchSeasonDTO) []interface{} {
	var results = seasonRepository.Search(orgOwnerID, leagueID, searchDTO)
	return results
}
