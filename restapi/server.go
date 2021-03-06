package restapi

import (
	"github.com/Rate-My-Bistro/crawler/config"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"log"
	"strconv"
)

// See Declarative Comments Format: https://swaggo.github.io/swaggo.io/declarative_comments_format/general_api_info.html

// @title This is a cgm bistro menu crawler
// @version 1.0.0
// @description This is a cgm bistro menu crawler

// @contact.name Rouven Himmelstein
// @contact.email rouven.himmelstein@cgm.com

// @host localhost:7331

// configures and starts the http server that serves the rest api
func Serve() {
	router := setupRouter()
	startRouter(router)
}

// adds routes to the server
func setupRouter() *gin.Engine {
	router := gin.Default()

	addApiDocEndpoint(router)
	addJobsResource(router)

	return router
}

// Define all routes for this resource
func addJobsResource(router *gin.Engine) {
	group := router.Group("/jobs")
	{
		group.GET("", jobGet())
		group.GET("/:id", jobGetWithParameter())

		group.POST("", jobPost())
	}
}

// adds the swagger api endpoint
func addApiDocEndpoint(router *gin.Engine) {
	restApiPort := strconv.FormatUint(config.Get().RestApiPort, 10)
	swaggerApiDocLocation := config.Get().SwaggerApiDocLocation

	router.StaticFile("/swagger.json", swaggerApiDocLocation)
	url := ginSwagger.URL("http://localhost:" + restApiPort + "/swagger.json")
	router.GET("/api/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

// starts the http server
func startRouter(server *gin.Engine) {
	portAsString := strconv.FormatUint(config.Get().RestApiPort, 10)
	log.Println("serving at http://localhost:" + portAsString)

	// This method will block the calling goroutine indefinitely unless an error happens.
	err := server.Run(":" + portAsString)

	if err != nil {
		log.Fatal(err)
	}
}
