package controllers

import "github.com/kataras/iris/v12"

type OrganizationsController struct{}

func NewOrganizationController() OrganizationsController {
	return OrganizationsController{}
}

func (oc *OrganizationsController) FetchOrganizations(ctx iris.Context) {
	ctx.JSON(iris.Map{
		"jwt": "hello",
	})
}
