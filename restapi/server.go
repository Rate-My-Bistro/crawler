package restapi

import (
	"github.com/ansgarS/rate-my-bistro-crawler/config"
	"github.com/yarf-framework/yarf"
	"log"
)

func Serve() {
	y := yarf.New()

	addResourceEndpoints(y)

	log.Println("🔥 serving from " + config.Cfg.RestApiAddress + " 🔥")
	y.Start(config.Cfg.RestApiAddress)
}
