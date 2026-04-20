package services

import (
	domainservices "league-management/internal/organization_management/domain/services"
	teampkg "league-management/internal/organization_management/domain/team"
	userdomain "league-management/internal/user_management/domain"
)

type TeamService struct {
	teamRepository   teamStore
	organizationRepo organizationFinder
	userRepo         teamUserFinder
}

type teamStore interface {
	Save(*teampkg.Team) error
	Search(string, string, string) []interface{}
	FindById(string) (*teampkg.Team, error)
}

type teamUserFinder interface {
	FindById(string) (*userdomain.User, error)
}

func NewTeamService(teamRepository teamStore, organizationRepo organizationFinder, userRepo teamUserFinder) *TeamService {
	return &TeamService{
		teamRepository:   teamRepository,
		organizationRepo: organizationRepo,
		userRepo:         userRepo,
	}
}

func (ts *TeamService) Make(teamName teampkg.TeamName, userId string, organizationId string) error {
	user, err := ts.userRepo.FindById(userId)

	if err != nil {
		return err
	}

	organization, err := ts.organizationRepo.FindById(organizationId)

	if err != nil {
		return err
	}

	newTeam, err := domainservices.CreateTeam(teamName, organization, user.Id)

	if err != nil {
		return err
	}

	err = ts.teamRepository.Save(newTeam)

	if err != nil {
		return err
	}

	return nil
}

func (ts *TeamService) Search(organizationId, userId, searchTerm string) []interface{} {

	result := ts.teamRepository.Search(organizationId, userId, searchTerm)

	return result
}
