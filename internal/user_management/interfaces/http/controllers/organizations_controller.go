package controllers

import (
	"github.com/kataras/iris/v12"
	"league-management/internal/organization_management/application/services"
	"league-management/internal/organization_management/domain"
	domain2 "league-management/internal/user_management/domain/user"
	"log"
)

var organizationService = services.NewOrganizationService()

type OrganizationsController struct{}

func NewOrganizationController() OrganizationsController {
	return OrganizationsController{}
}

type organizationRequestDto struct {
	Name string `json:"name"`
}

func (oc *OrganizationsController) FetchOrganizations(ctx iris.Context) {
	log.Println("organizations?")
	result := []domain.Organization{}
	ctx.JSON(result)
}

func (oc *OrganizationsController) AddOrganization(ctx iris.Context) {
	var createOrganizationRequestDto organizationRequestDto

	err := ctx.ReadJSON(&createOrganizationRequestDto)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	value := ctx.Values().Get("user")
	authenticatedUser, ok := value.(*domain2.User)

	log.Print(value, authenticatedUser, ok)

	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "No Authentication"})
		return
	}

	err = organizationService.OpenNewOrganization(authenticatedUser.Id, createOrganizationRequestDto.Name)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.JSON(iris.Map{"status": "ok"})
}
