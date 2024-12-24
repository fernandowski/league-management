package controllers

import (
	"github.com/kataras/iris/v12"
	"league-management/internal/organization_management/application/services"
	"league-management/internal/shared/dtos"
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

func (sc *SeasonController) Search(ctx iris.Context) {

	searchTerm := ctx.URLParamDefault("term", "")

	leagueId := ctx.Params().GetDefault("league_id", "").(string)
	if leagueId == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "missing league_id"})
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

	value := ctx.Values().Get("user")
	authenticatedUser, ok := value.(*domain2.User)

	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "No Authentication"})
		return
	}

	var searchDTO = dtos.SearchSeasonDTO{
		LeagueId: leagueId,
		BaseSearchDTO: dtos.BaseSearchDTO{
			Term:   searchTerm,
			Limit:  limit,
			Offset: offset,
		},
	}

	var results = seasonService.Search(authenticatedUser.Id, leagueId, searchDTO)

	ctx.JSON(results)
}

func (sc *SeasonController) Schedule(ctx iris.Context) {

	seasonId := ctx.Params().GetDefault("season_id", "").(string)
	if seasonId == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "missing season_id"})
		return
	}

	value := ctx.Values().Get("user")
	authenticatedUser, ok := value.(*domain2.User)

	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "No Authentication"})
		return
	}

	var err = seasonService.PlanSchedule(authenticatedUser.Id, seasonId)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.JSON(iris.Map{"status": "ok"})
}

func (sc *SeasonController) ChangeMatchScore(ctx iris.Context) {
	var body dtos.ChangeGameScoreDTO

	err := ctx.ReadJSON(&body)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	seasonId := ctx.Params().GetDefault("season_id", "").(string)
	if seasonId == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "missing season_id"})
		return
	}

	value := ctx.Values().Get("user")
	authenticatedUser, ok := value.(*domain2.User)

	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "No Authentication"})
		return
	}

	err = seasonService.ChangeMatchUpScore(authenticatedUser.Id, seasonId, body)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.JSON(iris.Map{"status": "ok"})

}

func (sc *SeasonController) StartSeason(ctx iris.Context) {

	seasonId := ctx.Params().GetDefault("season_id", "").(string)
	if seasonId == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "missing season_id"})
		return
	}

	value := ctx.Values().Get("user")
	authenticatedUser, ok := value.(*domain2.User)

	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "No Authentication"})
		return
	}

	var err = seasonService.StartSeason(authenticatedUser.Id, seasonId)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.JSON(iris.Map{"status": "ok"})
}

func (sc *SeasonController) SeasonDetails(ctx iris.Context) {

	seasonId := ctx.Params().GetDefault("season_id", "").(string)
	if seasonId == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "missing season_id"})
		return
	}

	value := ctx.Values().Get("user")
	authenticatedUser, ok := value.(*domain2.User)

	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "No Authentication"})
		return
	}

	result, err := seasonService.SeasonDetails(authenticatedUser.Id, seasonId)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.JSON(result)
}

func (sc *SeasonController) SeasonStandings(ctx iris.Context) {

	seasonId := ctx.Params().GetDefault("season_id", "").(string)
	if seasonId == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "missing season_id"})
		return
	}

	value := ctx.Values().Get("user")
	authenticatedUser, ok := value.(*domain2.User)

	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "No Authentication"})
		return
	}

	result, err := seasonService.SeasonStandings(authenticatedUser.Id, seasonId)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.JSON(result)
}
