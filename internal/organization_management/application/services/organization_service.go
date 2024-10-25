package services

import (
	"errors"
	"league-management/internal/organization_management/domain"
	"league-management/internal/organization_management/infrastructure/repositories"
	pg "league-management/internal/user_management/infrastructure/repositories/postgres"
	"log"
)

type OrganizationService struct {
}

var organizationRepo = repositories.NewOrganizationRepository()
var userRepo = pg.NewUserRepository()

func NewOrganizationService() *OrganizationService {
	return &OrganizationService{}
}

func (os *OrganizationService) FetchOrganizations(ownerId string) ([]domain.Organization, error) {

	// this needs to be an adapter
	// we are using the repo from another bounded context
	// in case we attempt to migrate to microservices having adapters
	// will make the migration easier for now it is okay.
	user, _ := userRepo.FindById(ownerId)

	if user == nil {
		return nil, errors.New("user does not exist")
	}

	organizations, err := organizationRepo.FetchAll(ownerId)
	log.Print(organizations)

	if err != nil {
		return nil, err
	}

	return organizations, nil
}

func (os *OrganizationService) OpenNewOrganization(organizationOwner string, orgName string) error {

	existingOrganization, _ := organizationRepo.FindByName(orgName, organizationOwner)

	if existingOrganization != nil {
		return errors.New("organization already exists with name")
	}

	organization := domain.NewOrganization(orgName, organizationOwner)

	err := organizationRepo.Save(&organization)
	if err != nil {
		return err
	}

	return nil
}
