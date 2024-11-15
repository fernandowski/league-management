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

func (ls *LeagueService) InitiateTeamMembership(orgOwnerID string, leagueID string, teamID string) error {

	organizationOwner, _ := userRepo.FindById(orgOwnerID)
	if organizationOwner == nil {
		return errors.New("user does not exist")
	}

	league, err := leagueRepository.FindById(leagueID)
	if err != nil {
		return errors.New("league does not exist")
	}

	organization, err := organizationRepo.FindById(league.OrganizationId)
	if belongsToUser := organization.BelongsToOwner(orgOwnerID); !belongsToUser {
		return errors.New("only organization owner allowed to change memberships")
	}

	team, err := teamRepository.FindById(teamID)
	if err != nil {
		return errors.New("team does not exist")
	}

	_, err = league.StartTeamMembership(string(*team.ID))
	if err != nil {
		return err
	}

	err = leagueRepository.Save(league)
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
