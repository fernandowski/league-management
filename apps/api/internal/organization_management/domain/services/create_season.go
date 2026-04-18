package services

import (
	"errors"
	"league-management/internal/organization_management/domain/league"
	"league-management/internal/organization_management/domain/organization"
	"league-management/internal/organization_management/domain/season"
)

func CreateSeason(
	organization *organization.Organization,
	league *league.League,
	orgOwnerID,
	seasonName string) (*season.Season, error) {

	if !organization.BelongsToOwner(orgOwnerID) {
		return nil, errors.New("only organization owner allowed to Add a new season")
	}

	if league.HasActiveSeason() {
		return nil, errors.New("league already has an active season")
	}

	newSeason, err := season.NewSeason(seasonName, *league.Id)

	if err != nil {
		return nil, err
	}

	return newSeason, nil
}
