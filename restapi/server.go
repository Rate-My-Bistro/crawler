package restapi

import (
	"github.com/ansgarS/rate-my-bistro-crawler/config"
	"github.com/yarf-framework/yarf"
	"log"
)

// starts the http server that serves the rest api
// the Start() command holds the application main thread
func Serve() {
	server := yarf.New()

	addJobsResource(server)

	log.Println("🔥 serving from " + config.Get().RestApiAddress + " 🔥")
	server.Start(config.Get().RestApiAddress)
}
