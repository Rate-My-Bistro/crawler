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
	DatabaseUser              string `env:"DATABASE_USER"`
	DatabasePassword          string `env:"DATABASE_PASSWORD"`
	MealCollectionName        string `env:"MEAL_COLLECTION_NAME"`
	JobCollectionName         string `env:"JOB_COLLECTION_NAME"`
	JobSchedulerTickInSeconds uint64 `env:"JOB_SCHEDULER_TICK_IN_SECONDS"`
	RestApiAddress            string `env:"REST_API_ADDRESS"`
}

var cfg Config
var cfgLoaded bool

// reads the application configuration from env files
// leave envPath blank to read the .env file from the current directory
func Get(envPaths ...string) Config {
	if cfgLoaded {
		return cfg
	}

	err := godotenv.Load(envPaths...)

	if err != nil {
		// if i am in a feature directory use the .env file from the parent directory
		err := godotenv.Load("../.env")

		if err != nil {
			log.Fatal("Error loading .env file "+
				" file not found, permission issues or a corrupted file",
				err)
		}
	}

	err = env.Parse(&cfg)

	if err != nil {
		log.Fatal("The .env format is not valid", err)
	}

	cfgLoaded = true

	return cfg
}
