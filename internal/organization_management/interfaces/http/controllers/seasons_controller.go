package controllers

import (
	"github.com/kataras/iris/v12"
	"league-management/internal/organization_management/application/services"
	domain2 "league-management/internal/user_management/domain/user"
)

type SeasonController struct {
}

var seasonService = services.NewSeasonService()

func NewSeasonController() *SeasonController {
	return &SeasonController{}
}
func (sc *SeasonController) AddNewSeasonToLeague(ctx iris.Context) {
	var requestBody struct {
		Name string `json:"name"`
	}

	if err := ctx.ReadJSON(&requestBody); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

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

	err := seasonService.AddNewSeason(authenticatedUser.Id, leagueId.(string), requestBody.Name)

	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.JSON(iris.Map{"status": "ok"})
}
