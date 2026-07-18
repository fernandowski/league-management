package controllers

import (
	"github.com/kataras/iris/v12"
	"league-management/internal/access_management/application/services"
	userdomain "league-management/internal/user_management/domain"
)

type AccessController struct {
	accessService *services.AccessService
}

type grantRequestDTO struct {
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	ScopeType string `json:"scope_type"`
	ScopeID   string `json:"scope_id"`
}

type grantResponseDTO struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	ScopeType string `json:"scope_type"`
	ScopeID   string `json:"scope_id"`
}

func NewAccessController(accessService *services.AccessService) *AccessController {
	return &AccessController{accessService: accessService}
}

func (ac *AccessController) GrantRole(ctx iris.Context) {
	var body grantRequestDTO
	if err := ctx.ReadJSON(&body); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	authenticatedUser, ok := authenticatedUser(ctx)
	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "No Authentication"})
		return
	}

	if err := ac.accessService.GrantRole(authenticatedUser.Id, body.UserID, body.Role, body.ScopeType, body.ScopeID); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.JSON(iris.Map{"status": "ok"})
}

func (ac *AccessController) RevokeRole(ctx iris.Context) {
	var body grantRequestDTO
	if err := ctx.ReadJSON(&body); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	authenticatedUser, ok := authenticatedUser(ctx)
	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "No Authentication"})
		return
	}

	if err := ac.accessService.RevokeRole(authenticatedUser.Id, body.UserID, body.Role, body.ScopeType, body.ScopeID); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.JSON(iris.Map{"status": "ok"})
}

func (ac *AccessController) ListGrants(ctx iris.Context) {
	userID := ctx.URLParamDefault("user_id", "")
	scopeType := ctx.URLParamDefault("scope_type", "")
	scopeID := ctx.URLParamDefault("scope_id", "")

	authenticatedUser, ok := authenticatedUser(ctx)
	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "No Authentication"})
		return
	}

	grants, err := ac.accessService.ListGrantsForActor(authenticatedUser.Id, userID, scopeType, scopeID)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	response := make([]grantResponseDTO, len(grants))
	for i, grant := range grants {
		response[i] = grantResponseDTO{
			ID:        grant.ID,
			UserID:    grant.UserID,
			Role:      string(grant.Role),
			ScopeType: string(grant.Scope.Type),
			ScopeID:   grant.Scope.ID,
		}
	}

	ctx.JSON(response)
}

func (ac *AccessController) CheckRole(ctx iris.Context) {
	userID := ctx.URLParamDefault("user_id", "")
	role := ctx.URLParamDefault("role", "")
	scopeType := ctx.URLParamDefault("scope_type", "")
	scopeID := ctx.URLParamDefault("scope_id", "")

	authenticatedUser, ok := authenticatedUser(ctx)
	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "No Authentication"})
		return
	}

	hasRole, err := ac.accessService.CheckRole(authenticatedUser.Id, userID, role, scopeType, scopeID)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.JSON(iris.Map{"has_role": hasRole})
}

func authenticatedUser(ctx iris.Context) (*userdomain.User, bool) {
	value := ctx.Values().Get("user")
	user, ok := value.(*userdomain.User)
	return user, ok
}
