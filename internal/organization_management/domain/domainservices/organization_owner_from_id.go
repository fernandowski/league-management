package domainservices

import (
	"errors"
	"league-management/internal/organization_management/domain"
)

func OrganizationOwnerFromUserId(userId *string, organization *domain.Organization) (*string, error) {
	if !organization.BelongsToOwner(*userId) {
		return nil, errors.New("organization does not belong to user")
	}

	return userId, nil
}
