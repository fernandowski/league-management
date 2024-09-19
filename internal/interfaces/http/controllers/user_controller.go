package controllers

import (
	"github.com/kataras/iris/v12"
	"league-management/internal/application/service"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(us *service.UserService) UserController {
	return UserController{userService: us}
}

func (uc *UserController) Register(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"message": "fernando"})
}
