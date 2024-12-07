package controllers

import (
	"github.com/kataras/iris/v12"
	"league-management/internal/organization_management/application/services"
	"league-management/internal/organization_management/domain"
	"league-management/internal/shared/dtos"
	domain2 "league-management/internal/user_management/domain/user"
)

var leagueService = services.NewLeagueService()

type LeaguesController struct{}

type leagueCreateDTO struct {
	Name           string `json:"name"`
	OrganizationId string `json:"organization_id"`
}

type leagueInvitationDTO struct {
	TeamId string `json:"team_id"`
}

type leagueSearchQueryDTO struct {
	OrganizationId string `json:"organization_id"`
}

type leagueResponseDto struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	TeamIds []string `json:"team_ids"`
}

func toLeagueSearchQueryDto(league domain.League) leagueResponseDto {
	return leagueResponseDto{
		Id:   *league.Id,
		Name: league.Name,
	}
}

func leaguesToRequestResponse(leagues []domain.League) []leagueResponseDto {
	dto := make([]leagueResponseDto, len(leagues))

	for i, league := range leagues {
		dto[i] = toLeagueSearchQueryDto(league)
	}
	return dto
}

func NewLeaguesController() *LeaguesController {
	return &LeaguesController{}
}
func (lc *LeaguesController) CreateLeague(ctx iris.Context) {
	var body leagueCreateDTO

	// TODO: THIS NEEDS TO BE A FUNCTION AT THIS POINT.
	err := ctx.ReadJSON(&body)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	value := ctx.Values().Get("user")
	authenticatedUser, ok := value.(*domain2.User)

	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "No Authentication"})
		return
	}

	err = leagueService.Provision(authenticatedUser.Id, body.OrganizationId, body.Name)

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.JSON(iris.Map{"status": "ok"})
}

func (lc *LeaguesController) StartLeagueMembership(ctx iris.Context) {
	var body leagueInvitationDTO
	err := ctx.ReadJSON(&body)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "missing request params"})
		return
	}

	value := ctx.Values().Get("user")
	authenticatedUser, ok := value.(*domain2.User)
	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "No Authentication"})
		return
	}

	teamId := ctx.Params().GetDefault("league_id", "")
	if teamId == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "missing league_id"})
		return
	}

	err = leagueService.InitiateTeamMembership(authenticatedUser.Id, teamId.(string), body.TeamId)

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.JSON(iris.Map{"status": "ok"})
}

func (lc *LeaguesController) RevokeLeagueMembership(ctx iris.Context) {
	leagueId := ctx.Params().GetDefault("league_id", "")
	if leagueId == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "missing league_id"})
		return
	}

	membershipId := ctx.Params().GetDefault("membership_id", "")
	if membershipId == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "missing membership_id"})
		return
	}

	value := ctx.Values().Get("user")
	authenticatedUser, ok := value.(*domain2.User)
	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "No Authentication"})
		return
	}

	err := leagueService.RevokeTeamMembership(authenticatedUser.Id, leagueId.(string), membershipId.(string))

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.JSON(iris.Map{"status": "ok"})
}

func (lc *LeaguesController) FetchLeagues(ctx iris.Context) {
	searchTerm := ctx.URLParamDefault("term", "")

	organizationId := ctx.URLParamDefault("organization_id", "")

	if organizationId == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Missing organization id"})
		return
	}

	value := ctx.Values().Get("user")
	authenticatedUser, ok := value.(*domain2.User)

	// TODO: If repeated again will create abstraction.
	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "No Authentication"})
		return
	}

	limit, _ := ctx.URLParamInt("limit")
	if limit < 0 {
		limit = 50
	}

	offset, _ := ctx.URLParamInt("offset")
	if offset < 0 {
		offset = 0
	}

	var searchDTO = dtos.LeagueSearchDTO{
		OrganizationID: organizationId,
		BaseSearchDTO: dtos.BaseSearchDTO{
			Limit:  limit,
			Offset: offset,
			Term:   searchTerm,
		},
	}

	leagues, err := leagueService.Search(authenticatedUser.Id, searchDTO)

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.JSON(leaguesToRequestResponse(leagues))
}

func (lc *LeaguesController) FetchLeaguesMembers(ctx iris.Context) {
	leagueId := ctx.Params().GetDefault("league_id", "")
	if leagueId == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "missing league_id"})
		return
	}

	value := ctx.Values().Get("user")
	authenticatedUser, ok := value.(*domain2.User)

	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "No Authentication"})
		return
	}

	result, err := leagueService.FetchLeagueMembers(leagueId.(string), authenticatedUser.Id)

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.JSON(result)
}
