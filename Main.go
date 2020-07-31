package main

import (
	"github.com/ansgarS/rate-my-bistro-crawler/crawler"
	"github.com/ansgarS/rate-my-bistro-crawler/persister"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	BistroUrl       string `env:"BISTRO_URL"`
	DatabaseAddress string `env:"DATABASE_ADDRESS"`
	DatabaseName    string `env:"DATABASE_NAME"`
	CollectionName  string `env:"COLLECTION_NAME"`
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

// handles the application cycle
func main() {
	crawledMeals := crawler.CrawlCurrentWeek(Cfg.BistroUrl)

	persister.PersistMeals(
		Cfg.DatabaseAddress,
		Cfg.DatabaseName,
		Cfg.CollectionName,
		crawledMeals)
}
