package http

import (
	accessservices "league-management/internal/access_management/application/services"
	accesspg "league-management/internal/access_management/infrastructure/repositories/postgres"
	accessroutes "league-management/internal/access_management/interfaces/http"
	accesscontrollers "league-management/internal/access_management/interfaces/http/controllers"
	orgservices "league-management/internal/organization_management/application/services"
	"league-management/internal/organization_management/infrastructure/repositories"
	orgroutes "league-management/internal/organization_management/interfaces/http"
	orgcontrollers "league-management/internal/organization_management/interfaces/http/controllers"
	userservices "league-management/internal/user_management/application/services"
	pg "league-management/internal/user_management/infrastructure/repositories/postgres"
	userroutes "league-management/internal/user_management/interfaces/http"
	usercontrollers "league-management/internal/user_management/interfaces/http/controllers"

	"github.com/kataras/iris/v12"
)

func NewApp() *iris.Application {
	app := iris.New()
	app.UseRouter(cors)
	app.Use(iris.Compression)

	userRepository := pg.NewUserRepository()
	organizationRepository := repositories.NewOrganizationRepository()
	teamRepository := repositories.NewTeamRepository()
	leagueRepository := repositories.NewLeagueRepository()
	seasonRepository := repositories.NewSeasonRepository()
	accessRepository := accesspg.NewAccessRepository()

	authService := userservices.NewAuthService()
	userService := userservices.NewUserService(userRepository, authService)
	accessService := accessservices.NewAccessService(accessRepository)
	organizationService := orgservices.NewOrganizationService(organizationRepository, userRepository)
	teamService := orgservices.NewTeamService(teamRepository, organizationRepository, userRepository)
	leagueService := orgservices.NewLeagueService(leagueRepository, organizationRepository, teamRepository, userRepository)
	seasonService := orgservices.NewSeasonService(seasonRepository, leagueRepository, organizationRepository)
	refereeService := orgservices.NewRefereeService()

	authMiddleware := NewAuthorizationMiddleware(authService, userRepository)

	userController := usercontrollers.NewUserController(userService)
	accessController := accesscontrollers.NewAccessController(accessService)
	organizationController := orgcontrollers.NewOrganizationController(organizationService)
	teamsController := orgcontrollers.NewTeamsController(teamService)
	leaguesController := orgcontrollers.NewLeaguesController(leagueService)
	seasonController := orgcontrollers.NewSeasonController(seasonService)
	refereeController := orgcontrollers.NewRefereeController(refereeService, seasonRepository)

	v1Router := app.Party("/v1")
	userroutes.RegisterRoutes(v1Router, userController)
	accessroutes.RegisterRoutes(v1Router, authMiddleware.Handler, accessController)
	orgroutes.RegisterRoutes(v1Router, authMiddleware.Handler, organizationController, teamsController, leaguesController, seasonController, refereeController)

	return app
}

func Run() error {
	return NewApp().Listen(":8080")
}

func cors(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	if ctx.Method() == iris.MethodOptions {
		ctx.Header("Access-Control-Allow-Methods", "POST, PUT, PATCH, DELETE, GET, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "*")
		ctx.Header("Access-Control-Max-Age", "86400")
		ctx.StatusCode(iris.StatusNoContent)
		return
	}

	ctx.Next()
}
