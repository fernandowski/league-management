package domainservices

import (
	"league-management/internal/organization_management/domain"
	"league-management/internal/organization_management/infrastructure/repositories"
)

type ProvisionLeagueService struct{}

func NewProvisionLeagueService() *ProvisionLeagueService {
	return &ProvisionLeagueService{}
}

var organizationRepository = repositories.NewOrganizationRepository()

var leagueRepository = repositories.NewLeagueRepository()

func (pls *ProvisionLeagueService) CanProvisionLeague(ownerId string, organizationId string, leagueName string) (*domain.League, error) {
	organization, err := organizationRepository.FindById(organizationId)

	if err != nil {
		return nil, err
	}

	if err := organization.CanAcceptANewLeague(ownerId); err != nil {
		return nil, err
	}

	newLeague := domain.NewLeague(leagueName, ownerId, organizationId)

	return &newLeague, nil
}
