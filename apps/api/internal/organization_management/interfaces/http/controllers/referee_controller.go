package controllers

import (
	"encoding/json"
	"league-management/internal/organization_management/application/services"
	"league-management/internal/organization_management/domain"
	"league-management/internal/organization_management/infrastructure/repositories"
	"league-management/internal/shared/dtos"
	"net/http"
)

type RefereeController struct {
	RefereeService   *services.RefereeService
	SeasonRepository *repositories.SeasonRepository
}

func NewRefereeController(refereeService *services.RefereeService, seasonRepo *repositories.SeasonRepository) *RefereeController {
	return &RefereeController{
		RefereeService:   refereeService,
		SeasonRepository: seasonRepo,
	}
}

func (rc *RefereeController) UpdateMatchScore(w http.ResponseWriter, r *http.Request) {
	var dto dtos.RefereeMatchUpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request body"))
		return
	}

	// Find the season containing the match
	season, err := rc.SeasonRepository.FindByID(dto.MatchID)
	if err != nil || season == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("season or match not found"))
		return
	}

	// Find the match in the season
	var matchToUpdate *domain.Match
	for i := range season.Rounds {
		for j := range season.Rounds[i].Matches {
			if season.Rounds[i].Matches[j].ID == dto.MatchID {
				matchToUpdate = &season.Rounds[i].Matches[j]
				break
			}
		}
	}

	if matchToUpdate == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("match not found"))
		return
	}

	err = rc.RefereeService.UpdateMatchScore(matchToUpdate, dto.RefereeID, dto.HomeScore, dto.AwayScore)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = rc.SeasonRepository.Save(season)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to save match update"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("match updated successfully"))
}
