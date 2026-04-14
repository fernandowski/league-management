package routes

import (
	"league-management/internal/organization_management/interfaces/http/controllers"

	"github.com/kataras/iris/v12"
)

func RegisterRoutes(
	router iris.Party,
	auth iris.Handler,
	organizationController *controllers.OrganizationsController,
	teamsController *controllers.TeamsController,
	leaguesController *controllers.LeaguesController,
	seasonController *controllers.SeasonController,
	refereeController *controllers.RefereeController,
) {
	organizationRouter := router.Party("/organizations", auth)
	{
		organizationRouter.Get("/", organizationController.FetchOrganizations)
		organizationRouter.Post("/", organizationController.AddOrganization)
	}

	leaguesRouter := router.Party("/leagues", auth)
	{
		leaguesRouter.Post("/", leaguesController.CreateLeague)
		leaguesRouter.Get("/", leaguesController.FetchLeagues)
		leaguesRouter.Get("/{league_id}", leaguesController.FetchLeagueDetails)
		leaguesRouter.Post("/{league_id}/invites", leaguesController.StartLeagueMembership)
		leaguesRouter.Get("/{league_id}/members", leaguesController.FetchLeaguesMembers)
		leaguesRouter.Delete("/{league_id}/members/{membership_id}", leaguesController.RevokeLeagueMembership)
		leaguesRouter.Post("/{league_id}/seasons", seasonController.AddNewSeasonToLeague)
		leaguesRouter.Get("/{league_id}/seasons", seasonController.Search)
		leaguesRouter.Post("/{league_id}/seasons/{season_id}/schedules", seasonController.Schedule)
		leaguesRouter.Get("/{league_id}/seasons/{season_id}", seasonController.SeasonDetails)
		leaguesRouter.Post("/{league_id}/seasons/{season_id}/start", seasonController.StartSeason)
		leaguesRouter.Put("/seasons/{season_id}/matches/score", seasonController.ChangeMatchScore)
	}

	teamsRouter := router.Party("/teams", auth)
	{
		teamsRouter.Post("/", teamsController.MakeTeam)
		teamsRouter.Get("/", teamsController.FetchAll)
	}

	seasonsRouter := router.Party("/seasons", auth)
	{
		seasonsRouter.Put("/{season_id}/matches/score", seasonController.ChangeMatchScore)
		seasonsRouter.Post("/{season_id}/rounds/complete", seasonController.CompleteCurrentRound)
		seasonsRouter.Get("/{season_id}/matches", seasonController.FetchMatches)
		seasonsRouter.Get("/{season_id}/standings", seasonController.SeasonStandings)
	}

	_ = refereeController
}
