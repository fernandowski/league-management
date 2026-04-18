package services

import (
	"errors"
	"league-management/internal/organization_management/domain/league"
	"league-management/internal/organization_management/domain/organization"
	"league-management/internal/organization_management/domain/team"
	"league-management/internal/shared/dtos"
	user_domain "league-management/internal/user_management/domain"
)

type leagueStore interface {
	FindById(string) (*league.League, error)
	Save(*league.League) error
	Search(dtos.LeagueSearchDTO) (*map[string]interface{}, error)
	FetchLeagueMembers(string) ([]interface{}, error)
	FetchLeagueDetails(league.League) (map[string]interface{}, error)
	WithLockedLeague(string, func(*league.League) error) error
}

type organizationFinder interface {
	FindById(string) (*organization.Organization, error)
}

type teamFinder interface {
	FindById(string) (*team.Team, error)
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

func NewLeagueService(
	leagueRepository leagueStore,
	organizationRepo organizationFinder,
	teamRepository teamFinder,
	userRepo userFinder,
) *LeagueService {
	return &LeagueService{
		leagueRepository: leagueRepository,
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

	organization, err := ls.organizationRepo.FindById(organizationId)
	if err != nil {
		return err
	}
	if err := organization.CanAcceptANewLeague(ownerId); err != nil {
		return err
	}

	newLeague := league.NewLeague("", leagueName, ownerId, organizationId)

	if err != nil {
		return err
	}

	err = ls.leagueRepository.Save(&newLeague)

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

	currentLeague, err := ls.leagueRepository.FindById(leagueID)
	if err != nil {
		return errors.New("league does not exist")
	}

	organization, err := ls.organizationRepo.FindById(currentLeague.OrganizationId)
	if belongsToUser := organization.BelongsToOwner(orgOwnerID); !belongsToUser {
		return errors.New("only organization owner allowed to change memberships")
	}

	team, err := ls.teamRepository.FindById(teamID)
	if err != nil {
		return errors.New("team does not exist")
	}

	return ls.leagueRepository.WithLockedLeague(leagueID, func(lockedLeague *league.League) error {
		_, err = lockedLeague.StartTeamMembership(string(*team.ID))
		return err
	})

}

func (ls *LeagueService) RevokeTeamMembership(orgOwnerID, leagueId, membershipId string) error {
	currentLeague, err := ls.leagueRepository.FindById(leagueId)
	if err != nil {
		return errors.New("league does not exist")
	}

	organization, err := ls.organizationRepo.FindById(currentLeague.OrganizationId)
	if belongsToUser := organization.BelongsToOwner(orgOwnerID); !belongsToUser {
		return errors.New("only organization owner allowed to change memberships")
	}

	return ls.leagueRepository.WithLockedLeague(leagueId, func(lockedLeague *league.League) error {
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
	currentLeague, err := ls.leagueRepository.FindById(leagueId)
	if err != nil {
		return nil, errors.New("league does not exist")
	}

	organization, err := ls.organizationRepo.FindById(currentLeague.OrganizationId)
	if err != nil {
		return nil, err
	}

	if belongs := organization.BelongsToOwner(userId); !belongs {
		return nil, errors.New("does not belong to user")
	}

	results, err := ls.leagueRepository.FetchLeagueMembers(*currentLeague.Id)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (ls *LeagueService) Details(userId, leagueId string) (map[string]interface{}, error) {
	currentLeague, err := ls.leagueRepository.FindById(leagueId)
	if err != nil {
		return nil, errors.New("league does not exist")
	}

	organization, err := ls.organizationRepo.FindById(currentLeague.OrganizationId)
	if err != nil {
		return nil, err
	}

	if belongs := organization.BelongsToOwner(userId); !belongs {
		return nil, errors.New("does not belong to user")
	}

	results, err := ls.leagueRepository.FetchLeagueDetails(*currentLeague)

	if err != nil {
		return nil, err
	}
	return results, nil
}
