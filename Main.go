package main

import (
	"fmt"
	"github.com/ansgarS/rate-my-bistro-crawler/persister"
	"log"
	"net/http"

	"github.com/ansgarS/rate-my-bistro-crawler/crawler"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	BistroUrl       string `env:"BISTRO_URL"`
	DatabaseAddress string `env:"DATABASE_ADDRESS"`
}

var Cfg Config

// the feature setup sequence
func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env.Parse(&Cfg)
}

func main() {
	httpResponse := readUrl(Cfg.BistroUrl)

	meals := crawler.Start(httpResponse.Body)

	persister.Start(Cfg.DatabaseAddress, meals)

	fmt.Print(meals)
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
