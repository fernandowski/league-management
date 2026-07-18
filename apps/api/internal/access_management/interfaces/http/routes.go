package http

import (
	"league-management/internal/access_management/interfaces/http/controllers"

	"github.com/kataras/iris/v12"
)

func RegisterRoutes(router iris.Party, auth iris.Handler, accessController *controllers.AccessController) {
	accessRouter := router.Party("/access", auth)
	{
		accessRouter.Post("/grants", accessController.GrantRole)
		accessRouter.Delete("/grants", accessController.RevokeRole)
		accessRouter.Get("/grants", accessController.ListGrants)
		accessRouter.Get("/check", accessController.CheckRole)
	}
}
