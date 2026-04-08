package http

import (
	"league-management/internal/user_management/application/services"
	domain "league-management/internal/user_management/domain"

	"github.com/kataras/iris/v12"
)

type userFinder interface {
	FindById(string) (*domain.User, error)
}

type AuthorizationMiddleware struct {
	authService *services.AuthService
	userRepo    userFinder
}

func NewAuthorizationMiddleware(authService *services.AuthService, userRepo userFinder) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		authService: authService,
		userRepo:    userRepo,
	}
}

func (am *AuthorizationMiddleware) Handler(ctx iris.Context) {
	type authorizationHeader struct {
		JWT string `header:"auth, required"`
	}

	var headers authorizationHeader
	if err := ctx.ReadHeaders(&headers); err != nil {
		ctx.StopExecution()
		return
	}

	claims, err := am.authService.ValidateJWT(headers.JWT)
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	user, err := am.userRepo.FindById(claims["user_id"].(string))
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.Values().Set("user", user)
	ctx.Next()
}
