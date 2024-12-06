package controllers

import (
	"github.com/kataras/iris/v12"
	"league-management/internal/organization_management/application/services"
	"league-management/internal/organization_management/domain"
	domain2 "league-management/internal/user_management/domain/user"
)

type TeamsController struct {
}

func NewTeamsController() *TeamsController {
	return &TeamsController{}
}

var teamService = services.NewTeamService()

type createTeamDto struct {
	Name           string `json:"name"`
	OrganizationId string `json:"organization_id"`
}

func (tc *TeamsController) MakeTeam(ctx iris.Context) {
	var createDto createTeamDto

	value := ctx.Values().Get("user")
	authenticatedUser, ok := value.(*domain2.User)

	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "No Authentication"})
		return
	}

	if err := ctx.ReadJSON(&createDto); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	err := teamService.Make(domain.TeamName(createDto.Name), authenticatedUser.Id, createDto.OrganizationId)

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.JSON(iris.Map{"status": "ok"})
}

func (tc *TeamsController) FetchAll(ctx iris.Context) {
	organizationId := ctx.URLParamDefault("organization_id", "")
	searchTerm := ctx.URLParamDefault("term", "")

	if organizationId == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Missing organization id"})
		return
	}

	value := ctx.Values().Get("user")
	authenticatedUser, ok := value.(*domain2.User)

	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "No Authentication"})
		return
	}

	results := teamService.Search(organizationId, authenticatedUser.Id, searchTerm)

	ctx.JSON(results)
}
