package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../services"
)

const (
	paramCountryId = "countryId"
)

func GetCountryFromAPI(context *gin.Context) {
	countryID := context.Param(paramCountryId)

	response, err := services.GetCountry(countryID)
	if err != nil {
		context.JSON(err.Status, err)
		return
	}

	context.JSON(http.StatusOK, response)
}