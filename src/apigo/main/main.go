package main

import (
	"github.com/gin-gonic/gin"
	"../controllers"
	"log"
	"../utils"
)

const (
	port = ":8080"
)

var (
	router = gin.Default()
)

func main() {
	utils.CircuitBreakerInstance = utils.NewState()

	router.GET("/users/:userId", controllers.GetUserFromAPI)
	router.GET("/countries/:countryId", controllers.GetCountryFromAPI)
	router.GET("/sites/:siteId", controllers.GetSiteFromAPI)
	router.GET("/results/:userId", controllers.GetResultFromAPI)

	//ruta para result workgroup
	router.GET("/waitgroup/results/:userId", controllers.GetResultFromApiWithWaitGroup)

	//ruta para result con channels
	router.GET("/channel/results/:userId", controllers.GetResultFromApiWithChannel)

	log.Fatal(router.Run(port))
}
