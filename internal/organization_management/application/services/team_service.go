package services

import (
	"league-management/internal/organization_management/domain"
	"league-management/internal/organization_management/domain/domainservices"
	"league-management/internal/organization_management/infrastructure/repositories"
)

type TeamService struct {
}

func NewTeamService() *TeamService {
	return &TeamService{}
}

var teamRepository = repositories.NewTeamRepository()

func (ts *TeamService) Make(teamName domain.TeamName, userId string, organizationId string) error {
	user, err := userRepo.FindById(userId)

	if err != nil {
		return err
	}

	organization, err := organizationRepo.FindById(organizationId)

	if err != nil {
		return err
	}

	newTeam, err := domainservices.CreateTeam(teamName, organization, user.Id)

	if err != nil {
		return err
	}

	err = teamRepository.Save(newTeam)

	if err != nil {
		return err
	}

	return nil
}

func (ts *TeamService) Search(organizationId string, userId string) []interface{} {

	result := teamRepository.FetchAll(organizationId, userId)

	return result
}
