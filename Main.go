package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ansgarS/rate-my-bistro-crawler/crawler"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	BistroUrl string `env:"BISTRO_URL"`
}

func main() {
	cfg := initConfig()
	httpResponse := readUrl(cfg.BistroUrl)

	meals := crawler.Start(httpResponse.Body)

	fmt.Print(meals)
}

func initConfig() (cfg Config) {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env.Parse(&cfg)

	return cfg
}

func readUrl(url string) *http.Response {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	return res
}
