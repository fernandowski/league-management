package services

import (
	"errors"
	"league-management/internal/organization_management/domain"
	"league-management/internal/organization_management/domain/domainservices"
	"league-management/internal/organization_management/infrastructure/repositories"
)

type LeagueService struct{}

func NewLeagueService() *LeagueService {
	return &LeagueService{}
}

var leagueRepository = repositories.NewLeagueRepository()

func (ls *LeagueService) Provision(ownerId string, organizationId string, leagueName string) error {

	user, _ := userRepo.FindById(ownerId)

	if user == nil {
		return errors.New("user does not exist")
	}

	provisionLeagueService := domainservices.NewProvisionLeagueService()
	newLeague, err := provisionLeagueService.CanProvisionLeague(ownerId, organizationId, leagueName)

	if err != nil {
		return err
	}

	err = leagueRepository.Save(newLeague)

	if err != nil {
		return err
	}

	return nil
}

func (ls *LeagueService) FetchLeagues(userId string, organizationId string) ([]domain.League, error) {
	user, _ := userRepo.FindById(userId)

	if user == nil {
		return nil, errors.New("user does not exist")
	}

	organization, err := organizationRepo.FindById(organizationId)
	if err != nil {
		return nil, err
	}

	if belongs := organization.BelongsToOwner(user.Id); !belongs {
		return nil, errors.New("does not belong to user")
	}

	leagues, err := leagueRepository.FetchAll(organizationId)

	if err != nil {
		return nil, err
	}

	return leagues, nil
}
