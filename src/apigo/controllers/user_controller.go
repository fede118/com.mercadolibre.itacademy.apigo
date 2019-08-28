package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../services"
	"strconv"
	"../utils"
)

const (
	paramUserId = "userId"
)

func GetUserFromAPI(context *gin.Context) {
	userID, convError := strconv.Atoi(context.Param(paramUserId))
	if convError != nil {
		apiErr := utils.ApiError{
			Message: "error converting userID to INT",
			Status: http.StatusBadRequest,
		}
		context.JSON(apiErr.Status, apiErr)
		return
	}

	response, apiErr := services.GetUser(userID)
	if apiErr != nil {
		context.JSON(apiErr.Status, apiErr)
		return
	}

	context.JSON(http.StatusOK, response)
}