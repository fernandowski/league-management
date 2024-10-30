package services

import (
	"errors"
	"league-management/internal/organization_management/domain"
	"league-management/internal/organization_management/infrastructure/repositories"
)

type LeagueService struct{}

func NewLeagueService() *LeagueService {
	return &LeagueService{}
}

var leagueRepository = repositories.NewLeagueRepository()

func (ls *LeagueService) Provision(ownerId string, organizationId string, leagueName string) error {
	organization, err := organizationRepo.FindById(organizationId)

	user, _ := userRepo.FindById(ownerId)

	if user == nil {
		return errors.New("user does not exist")
	}

	if err != nil {
		return err
	}

	if !organization.IsInGoodStanding() {
		return errors.New("organization Not in good state")
	}

	if !organization.BelongsToOwner(ownerId) {
		return errors.New("cannot provision league")
	}

	newLeague := domain.NewLeague(leagueName, ownerId, organizationId)

	err = leagueRepository.Save(&newLeague)

	if err != nil {
		return err
	}

	return nil
}
