package controllers

import (
	"github.com/kataras/iris/v12"
	"league-management/internal/organization_management/application/services"
	"league-management/internal/organization_management/domain/organization"
	domain2 "league-management/internal/user_management/domain"
)

type OrganizationsController struct {
	organizationService *services.OrganizationService
}

func NewOrganizationController(organizationService *services.OrganizationService) *OrganizationsController {
	return &OrganizationsController{organizationService: organizationService}
}

type organizationRequestDto struct {
	Name string `json:"name"`
}

type organizationResponseDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func toOrganizationDto(organization organization.Organization) organizationResponseDto {
	return organizationResponseDto{
		Id:   *organization.ID,
		Name: organization.Name,
	}
}
func organizationToRequestResponse(organizations []organization.Organization) []organizationResponseDto {
	dto := make([]organizationResponseDto, len(organizations))

	for i, organization := range organizations {
		dto[i] = toOrganizationDto(organization)
	}

	return dto
}

func (oc *OrganizationsController) FetchOrganizations(ctx iris.Context) {

	value := ctx.Values().Get("user")
	authenticatedUser, ok := value.(*domain2.User)

	// TODO: If repeated again will create abstraction.
	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "No Authentication"})
		return
	}

	organizations, err := oc.organizationService.FetchOrganizations(authenticatedUser.Id)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	// TODO: We will need utility functions to always include data and status ok for now.
	ctx.JSON(organizationToRequestResponse(organizations))
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

	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "No Authentication"})
		return
	}

	err = oc.organizationService.OpenNewOrganization(authenticatedUser.Id, createOrganizationRequestDto.Name)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.JSON(iris.Map{"status": "ok"})
}
