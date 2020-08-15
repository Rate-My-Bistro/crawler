package restapi

import (
	"github.com/ansgarS/rate-my-bistro-crawler/config"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

// starts the http server that serves the rest api
// the Start() command block the calling goroutine
func Serve() {
	server := setupRouter()
	startRouter(server)
}

func setupRouter() *gin.Engine {
	server := gin.Default()

	addJobsResource(server)

	return server
}

func startRouter(server *gin.Engine) {
	portAsString := strconv.FormatUint(config.Get().RestApiPort, 10)

	// This method will block the calling goroutine indefinitely unless an error happens.
	err := server.Run(":" + portAsString)

	if err != nil {
		log.Fatal(err)
	}
}
