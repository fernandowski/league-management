package domainservices

import (
	"errors"
	"league-management/internal/organization_management/domain"
)

func CreateSeason(
	organization *domain.Organization,
	league *domain.League,
	orgOwnerID,
	seasonName string) (*domain.Season, error) {

	if !organization.BelongsToOwner(orgOwnerID) {
		return nil, errors.New("only organization owner allowed to Add a new season")
	}

	if league.HasActiveSeason() {
		return nil, errors.New("league already has an active season")
	}

	newSeason, err := domain.NewSeason(seasonName, *league.Id)

	if err != nil {
		return nil, err
	}

	return newSeason, nil
}
