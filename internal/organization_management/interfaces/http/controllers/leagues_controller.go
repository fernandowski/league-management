package controllers

import (
	"github.com/kataras/iris/v12"
	"league-management/internal/organization_management/application/services"
	domain2 "league-management/internal/user_management/domain/user"
	"log"
)

var leagueService = services.NewLeagueService()

type LeaguesController struct{}

type leagueCreateDTO struct {
	Name           string `json:"name"`
	OrganizationId string `json:"organization_id"`
}

type leagueSearchQueryDTO struct {
	OrganizationId string `json:"organization_id"`
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

func (lc *LeaguesController) FetchLeagues(ctx iris.Context) {
	organizationId := ctx.URLParamDefault("organization_id", "")
	log.Print(organizationId)

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

	leagues, err := leagueService.FetchLeagues(authenticatedUser.Id, organizationId)

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.JSON(leagues)
}