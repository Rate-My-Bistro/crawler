package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	BistroUrl                 string `env:"BISTRO_URL"`
	DatabaseAddress           string `env:"DATABASE_ADDRESS"`
	DatabaseName              string `env:"DATABASE_NAME"`
	MealCollectionName        string `env:"MEAL_COLLECTION_NAME"`
	JobCollectionName         string `env:"JOB_COLLECTION_NAME"`
	JobSchedulerTickInSeconds uint64 `env:"JOB_SCHEDULER_TICK_IN_SECONDS"`
}

var Cfg Config

// reads the application configuration from the env file.
func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file "+
			"due to permission issues or a corrupted file",
			err)
	}

	err = env.Parse(&Cfg)

	if err != nil {
		log.Fatal("The .env format is not valid", err)
	}
}
