package restapi

import (
	"github.com/ansgarS/rate-my-bistro-crawler/config"
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
	router.StaticFile("/docs/swagger.json", "./restapi/docs/swagger.json")
	url := ginSwagger.URL("http://localhost:" + strconv.FormatUint(config.Get().RestApiPort, 10) + "/docs/swagger.json")
	router.GET("/api/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

// starts the http server
func startRouter(server *gin.Engine) {
	portAsString := strconv.FormatUint(config.Get().RestApiPort, 10)

	// This method will block the calling goroutine indefinitely unless an error happens.
	log.Println("Serving at " + portAsString)
	err := server.Run(":" + portAsString)

	if err != nil {
		log.Fatal(err)
	}
}
