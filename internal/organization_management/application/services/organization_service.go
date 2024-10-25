package services

import (
	"errors"
	"league-management/internal/organization_management/domain"
	"league-management/internal/organization_management/infrastructure/repositories"
)

type OrganizationService struct {
}

var organizationRepo = repositories.NewOrganizationRepository()

func NewOrganizationService() *OrganizationService {
	return &OrganizationService{}
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
