package routes

import (
	"github.com/kataras/iris/v12"
	"league-management/internal/user_management/application/service"
	pg "league-management/internal/user_management/infrastructure/repositories/postgres"
	"league-management/internal/user_management/interfaces/http/controllers"
	"log"
)

func registerV1API(app *iris.Application) {
	v1Router := app.Party("/v1")
	initRoutes(v1Router)
}

func initRoutes(router iris.Party) {
	userService := service.NewUserService()
	userController := controllers.NewUserController(userService)

	initUserRouter(router, userController)
	initOrganizationRouter(router)
}

func initUserRouter(router iris.Party, uc controllers.UserController) {
	userRouter := router.Party("/user")
	{
		userRouter.Post("/register", uc.Register)
		userRouter.Post("/login", uc.Login)
	}
}

func initOrganizationRouter(router iris.Party) {

	var organizationController = controllers.NewOrganizationController()
	var organizationRouter = router.Party("/organizations", authorizationMiddleWare)
	{
		organizationRouter.Get("/", organizationController.FetchOrganizations)
		organizationRouter.Post("/", organizationController.AddOrganization)
	}
}

func loginEndpoint(ctx iris.Context) {
	ctx.WriteString("Login Endpoint")
}

func authorizationMiddleWare(ctx iris.Context) {
	type authorizationHeader struct {
		JWT string `header:"auth, required"`
	}

	var headers authorizationHeader
	if err := ctx.ReadHeaders(&headers); err != nil {
		ctx.StopExecution()
		return
	}

	var authorizationService = service.NewAuthService()

	claims, err := authorizationService.ValidateJWT(headers.JWT)

	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	var userRepository = pg.NewUserRepository()
	user, err := userRepository.FindById(claims["user_id"].(string))

	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.Values().Set("user", user)
	ctx.Next()
}

func Initialize() *iris.Application {
	app := iris.New()

	crs := func(ctx iris.Context) {
		ctx.Header("Access-Control-Allow-Origin", "http://localhost:8081")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Method() == iris.MethodOptions {
			ctx.Header("Access-Control-Methods", "POST, PUT, PATCH, DELETE, GET")
			ctx.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,Content-Type,auth")
			ctx.Header("Access-Control-Max-Age", "86400")
			ctx.StatusCode(iris.StatusNoContent)
			return
		}

		log.Print("here")

		ctx.Next()
	}

	app.UseRouter(crs)

	app.Use(iris.Compression)

	registerV1API(app)

	app.Listen(":8080")

	return app
}
