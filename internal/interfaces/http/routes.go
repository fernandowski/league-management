package routes

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"league-management/internal/application/service"
	"league-management/internal/interfaces/http/controllers"
)

func registerV1API(app *iris.Application) {
	v1Router := app.Party("/v1")
	initRoutes(v1Router)
}

func initRoutes(router iris.Party) {
	userService := service.NewUserService()
	userController := controllers.NewUserController(userService)
	initUserRouter(router, userController)
}

func initUserRouter(router iris.Party, uc controllers.UserController) {
	userRouter := router.Party("/user")
	{
		userRouter.Post("/register", uc.Register)
	}
}

func loginEndpoint(ctx iris.Context) {
	ctx.WriteString("Login Endpoint")
}

func Initialize() *iris.Application {
	app := iris.New()

	app.Use(iris.Compression)

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"},                   // Specify the origin explicitly
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Specify allowed methods
		AllowedHeaders:   []string{"Content-Type", "Authorization"},           // Specify allowed headers
		AllowCredentials: true,                                                // Allow credentials
	})
	app.UseRouter(crs)

	registerV1API(app)

	app.Listen(":8080")

	return app
}
