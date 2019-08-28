package controllers

import (
	"github.com/gin-gonic/gin"
	"../services"
	"net/http"
)

const (
	paramSiteId = "siteId"
)

func GetSiteFromAPI(context *gin.Context) {
	siteID := context.Param(paramSiteId)

	response, err := services.GetSite(siteID)
	if err != nil {
		context.JSON(err.Status, err)
		return
	}

	context.JSON(http.StatusOK, response)
}
