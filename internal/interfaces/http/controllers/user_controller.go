package controllers

import (
	"github.com/kataras/iris/v12"
	"league-management/internal/application/service"
	"league-management/internal/infrastructure/dto"
)

type UserController struct {
	userService *service.UserService
}

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserController(us *service.UserService) UserController {
	return UserController{userService: us}
}

func (uc *UserController) Register(ctx iris.Context) {

	var body userRequest
	err := ctx.ReadJSON(&body)

	createdUser, err := uc.userService.RegisterUser(body.Email, body.Password)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "bad request"})
		return
	}

	userDto := dto.FromDomain(createdUser)
	ctx.JSON(userDto)
}
