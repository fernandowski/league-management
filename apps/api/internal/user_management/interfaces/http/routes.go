package routes

import (
	"league-management/internal/user_management/interfaces/http/controllers"

	"github.com/kataras/iris/v12"
)

func RegisterRoutes(router iris.Party, userController *controllers.UserController) {
	userRouter := router.Party("/user")
	{
		userRouter.Post("/register", userController.Register)
		userRouter.Post("/login", userController.Login)
	}
}
