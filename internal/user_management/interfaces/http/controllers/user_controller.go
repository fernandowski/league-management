package controllers

import (
	"github.com/kataras/iris/v12"
	"league-management/internal/user_management/application/service"
	"league-management/internal/user_management/infrastructure/dto"
	"log"
)

type UserController struct {
	userService *service.UserService
}

type userRegistrationDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userLoginRequestDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserController(us *service.UserService) UserController {
	return UserController{userService: us}
}

func (uc *UserController) Register(ctx iris.Context) {

	var body userRegistrationDto
	err := ctx.ReadJSON(&body)

	createdUser, err := uc.userService.RegisterUser(body.Email, body.Password)
	if err != nil {
		log.Print(err)
		ctx.StatusCode(iris.StatusConflict)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	userDto := dto.FromDomain(createdUser)
	ctx.JSON(userDto)
}

func (uc *UserController) Login(ctx iris.Context) {
	var body userLoginRequestDto
	err := ctx.ReadJSON(&body)

	jwt, err := uc.userService.Login(body.Email, body.Password)
	if err != nil {
		log.Print(err)
		ctx.StatusCode(iris.StatusConflict)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.JSON(iris.Map{
		"jwt": jwt,
	})
}
