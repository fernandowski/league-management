package services

import (
	"errors"
	"league-management/internal/organization_management/domain"
	"league-management/internal/organization_management/infrastructure/repositories"
)

type SeasonService struct {
}

func NewSeasonService() *SeasonService {
	return &SeasonService{}
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
