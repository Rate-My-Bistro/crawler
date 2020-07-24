package main

import (
	"github.com/ansgarS/rate-my-bistro-crawler/persister"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

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

	err = env.Parse(&Cfg)

	if err != nil {
		log.Fatal("Error parsing .env file")
	}
}

func main() {
	documentReader := createBistroReader(Cfg.BistroUrl)

	crawledMeals := crawler.Start(documentReader)

	persister.Start(Cfg.DatabaseAddress, crawledMeals)
}

func createBistroReader(bistroUrl string) (documentReader io.Reader) {

	if strings.HasPrefix(bistroUrl, "file://") {
		bistroUrl := strings.Replace(bistroUrl, "file://", "", -1)
		documentReader = readFile(bistroUrl)
	} else {
		documentReader = readUrl(bistroUrl).Body
	}

	return documentReader
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

func readFile(filePath string) *os.File {
	bistroPageReader, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	return bistroPageReader
}
