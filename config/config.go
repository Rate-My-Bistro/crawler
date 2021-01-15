package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
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
	RestApiPort               uint64 `env:"REST_API_PORT"`
	SwaggerApiDocLocation     string `env:"SWAGGER_API_DOC_LOCATION"`
}

var cfg Config
var cfgPresent bool

// reads the application configuration from env files
// leave envPath blank to read the .env file from the current directory
func Get() Config {
	if cfgPresent {
		return cfg
	}

	if isTesting() {
		loadTestConfig()
	} else {
		loadConfig()
	}

	cfgPresent = true

	return cfg
}

// Tests if the application runs in a test context
func isTesting() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test") {
			return true
		}
	}
	return false
}

// loads the configuration from the .env.testing file
func loadTestConfig() {
	err := godotenv.Load("../.env.testing")

	if err != nil {
		log.Fatal("Error loading .env.testing file "+
			"not found, permission issues or a corrupted file",
			err)
	}

	err = env.Parse(&cfg)

	if err != nil {
		log.Fatal("The format is not valid", err)
	}
}

// load the configuration from the specified env path
// or dynamically discovers the .env file location
func loadConfig() {
	err := godotenv.Load()

	if err != nil {
		log.Println("did not load any .env file")
	}

	err = env.Parse(&cfg)

	if err != nil {
		log.Fatal("The .env format is not valid", err)
	}
}
