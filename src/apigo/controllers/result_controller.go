package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../services"
	"strconv"
	"../utils"
	)

func GetResultFromAPI(context *gin.Context) {
	userID, convError := strconv.Atoi(context.Param(paramUserId))
	if convError != nil {
		apiErr := utils.ApiError{
			Message: "error converting userID to INT",
			Status: http.StatusBadRequest,
		}
		context.JSON(apiErr.Status, apiErr)
		return
	}


	response, apiErr := services.GetResult(userID)
	if apiErr != nil {
		context.JSON(apiErr.Status, apiErr)
		return
	}

	context.JSON(http.StatusOK, response)
}

func GetResultFromApiWithWaitGroup(context *gin.Context) {
	userID, convError := strconv.Atoi(context.Param(paramUserId))
	if convError != nil {
		apiErr := utils.ApiError{
			Message: "error converting userID to INT",
			Status: http.StatusBadRequest,
		}
		context.JSON(apiErr.Status, apiErr)
		return
	}


	response, apiErr := services.GetResultWithWaitGroup(userID)
	if apiErr != nil {
		context.JSON(apiErr.Status, apiErr)
		return
	}

	context.JSON(http.StatusOK, response)
}

func GetResultFromApiWithChannel(context *gin.Context) {
	if utils.CircuitBreakerInstance.State != utils.CLOSED {
		println("CircuitBreaker is not CLOSED")
		apiErr := utils.ApiError{
			Message: "CIRCUITBREAKER ==> Connection is not OPEN",
			Status: http.StatusInternalServerError,
		}
		context.JSON(apiErr.Status, apiErr)
		return
	}


	userID, convError := strconv.Atoi(context.Param(paramUserId))
	if convError != nil {
		apiErr := utils.ApiError{
			Message: "error converting userID to INT",
			Status: http.StatusBadRequest,
		}
		context.JSON(apiErr.Status, apiErr)
		return
	}


	response, apiErr := services.GetResultWithChannel(userID)
	if apiErr != nil {
		context.JSON(apiErr.Status, apiErr)
		return
	}

	utils.CircuitBreakerInstance.Reset()
	context.JSON(http.StatusOK, response)
}
