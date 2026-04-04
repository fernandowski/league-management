package services

import (
	"errors"
	"league-management/internal/organization_management/domain"
	"league-management/internal/organization_management/domain/domainservices"
	"league-management/internal/organization_management/infrastructure/repositories"
	"league-management/internal/shared/dtos"
	user_domain "league-management/internal/user_management/domain/user"
)

type leagueStore interface {
	FindById(string) (*domain.League, error)
	Save(*domain.League) error
	Search(dtos.LeagueSearchDTO) (*map[string]interface{}, error)
	FetchLeagueMembers(string) ([]interface{}, error)
	FetchLeagueDetails(domain.League) (map[string]interface{}, error)
	WithLockedLeague(string, func(*domain.League) error) error
}

type organizationFinder interface {
	FindById(string) (*domain.Organization, error)
}

type teamFinder interface {
	FindById(string) (*domain.Team, error)
}

type userFinder interface {
	FindById(string) (*user_domain.User, error)
}

type LeagueService struct {
	leagueRepository leagueStore
	organizationRepo organizationFinder
	teamRepository   teamFinder
	userRepo         userFinder
}

func NewLeagueService() *LeagueService {
	return &LeagueService{
		leagueRepository: repositories.NewLeagueRepository(),
		organizationRepo: organizationRepo,
		teamRepository:   teamRepository,
		userRepo:         userRepo,
	}
}

func (ls *LeagueService) Provision(ownerId string, organizationId string, leagueName string) error {

	user, _ := ls.userRepo.FindById(ownerId)

	if user == nil {
		return errors.New("user does not exist")
	}

	provisionLeagueService := domainservices.NewProvisionLeagueService()
	newLeague, err := provisionLeagueService.CanProvisionLeague(ownerId, organizationId, leagueName)

	if err != nil {
		return err
	}

	err = ls.leagueRepository.Save(newLeague)

	if err != nil {
		return err
	}

	return nil
}

func (ls *LeagueService) InitiateTeamMembership(orgOwnerID string, leagueID string, teamID string) error {

	organizationOwner, _ := ls.userRepo.FindById(orgOwnerID)
	if organizationOwner == nil {
		return errors.New("user does not exist")
	}

	league, err := ls.leagueRepository.FindById(leagueID)
	if err != nil {
		return errors.New("league does not exist")
	}

	organization, err := ls.organizationRepo.FindById(league.OrganizationId)
	if belongsToUser := organization.BelongsToOwner(orgOwnerID); !belongsToUser {
		return errors.New("only organization owner allowed to change memberships")
	}

	team, err := ls.teamRepository.FindById(teamID)
	if err != nil {
		return errors.New("team does not exist")
	}

	return ls.leagueRepository.WithLockedLeague(leagueID, func(lockedLeague *domain.League) error {
		_, err = lockedLeague.StartTeamMembership(string(*team.ID))
		return err
	})

}

func (ls *LeagueService) RevokeTeamMembership(orgOwnerID, leagueId, membershipId string) error {
	league, err := ls.leagueRepository.FindById(leagueId)
	if err != nil {
		return errors.New("league does not exist")
	}

	organization, err := ls.organizationRepo.FindById(league.OrganizationId)
	if belongsToUser := organization.BelongsToOwner(orgOwnerID); !belongsToUser {
		return errors.New("only organization owner allowed to change memberships")
	}

	return ls.leagueRepository.WithLockedLeague(leagueId, func(lockedLeague *domain.League) error {
		_, err = lockedLeague.RemoveMembership(membershipId)
		return err
	})
}

func (ls *LeagueService) Search(userId string, searchDTO dtos.LeagueSearchDTO) (*map[string]interface{}, error) {
	user, _ := ls.userRepo.FindById(userId)

	if user == nil {
		return nil, errors.New("user does not exist")
	}

	organization, err := ls.organizationRepo.FindById(searchDTO.OrganizationID)
	if err != nil {
		return nil, err
	}

	if belongs := organization.BelongsToOwner(user.Id); !belongs {
		return nil, errors.New("does not belong to user")
	}

	leagues, err := ls.leagueRepository.Search(searchDTO)

	if err != nil {
		return nil, err
	}

	return leagues, nil
}

func (ls *LeagueService) FetchLeagueMembers(leagueId, userId string) ([]interface{}, error) {
	league, err := ls.leagueRepository.FindById(leagueId)
	if err != nil {
		return nil, errors.New("league does not exist")
	}

	organization, err := ls.organizationRepo.FindById(league.OrganizationId)
	if err != nil {
		return nil, err
	}

	if belongs := organization.BelongsToOwner(userId); !belongs {
		return nil, errors.New("does not belong to user")
	}

	results, err := ls.leagueRepository.FetchLeagueMembers(*league.Id)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (ls *LeagueService) Details(userId, leagueId string) (map[string]interface{}, error) {
	league, err := ls.leagueRepository.FindById(leagueId)
	if err != nil {
		return nil, errors.New("league does not exist")
	}

	organization, err := ls.organizationRepo.FindById(league.OrganizationId)
	if err != nil {
		return nil, err
	}

	if belongs := organization.BelongsToOwner(userId); !belongs {
		return nil, errors.New("does not belong to user")
	}

	results, err := ls.leagueRepository.FetchLeagueDetails(*league)

	if err != nil {
		return nil, err
	}
	return results, nil
}
