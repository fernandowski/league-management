package services

import (
	"errors"
	organizationpkg "league-management/internal/organization_management/domain/organization"
	userdomain "league-management/internal/user_management/domain"
)

type OrganizationService struct {
	organizationRepo organizationRepository
	userRepo         organizationUserFinder
}

type organizationRepository interface {
	FindByName(string, string) (*organizationpkg.Organization, error)
	FindById(string) (*organizationpkg.Organization, error)
	Save(*organizationpkg.Organization) error
	FetchAll(string) ([]organizationpkg.Organization, error)
}

type organizationUserFinder interface {
	FindById(string) (*userdomain.User, error)
}

func NewOrganizationService(organizationRepo organizationRepository, userRepo organizationUserFinder) *OrganizationService {
	return &OrganizationService{
		organizationRepo: organizationRepo,
		userRepo:         userRepo,
	}
}

func (os *OrganizationService) FetchOrganizations(ownerId string) ([]organizationpkg.Organization, error) {

	// this needs to be an adapter
	// we are using the repo from another bounded context
	// in case we attempt to migrate to microservices having adapters
	// will make the migration easier for now it is okay.
	user, _ := os.userRepo.FindById(ownerId)

	if user == nil {
		return nil, errors.New("user does not exist")
	}

	organizations, err := os.organizationRepo.FetchAll(ownerId)

	if err != nil {
		return nil, err
	}

	return organizations, nil
}

func (os *OrganizationService) OpenNewOrganization(organizationOwner string, orgName string) error {

	existingOrganization, _ := os.organizationRepo.FindByName(orgName, organizationOwner)

	if existingOrganization != nil {
		return errors.New("organization already exists with name")
	}

	newOrganization := organizationpkg.NewOrganization(nil, orgName, organizationOwner, true)

	err := os.organizationRepo.Save(&newOrganization)
	if err != nil {
		return err
	}

	return nil
}
