package routes

import (
	"github.com/kataras/iris/v12"
	"league-management/internal/user_management/application/service"
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

	crs := func(ctx iris.Context) {
		ctx.Header("Access-Control-Allow-Origin", "http://localhost:8081")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Method() == iris.MethodOptions {
			ctx.Header("Access-Control-Methods", "POST, PUT, PATCH, DELETE")
			ctx.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,Content-Type")
			ctx.Header("Access-Control-Max-Age", "86400")
			ctx.StatusCode(iris.StatusNoContent)
			return
		}

		log.Print("here")

		ctx.Next()
	} // or	"github.com/iris-contrib/middleware/cors"

	app.UseRouter(crs)

	app.Use(iris.Compression)

	registerV1API(app)

	app.Listen(":8080")

	return app
}
