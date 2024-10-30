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

	log.Print(ok)

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
