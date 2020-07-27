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
	documentReader := createBistroReader(Cfg.BistroUrl)

	crawledMeals := crawler.Crawl(documentReader)

	persister.PersistMeals(
		Cfg.DatabaseAddress,
		Cfg.DatabaseName,
		Cfg.CollectionName,
		crawledMeals)
}

// creates an reader object based on the provided bistroUrl
func createBistroReader(bistroUrl string) (documentReader io.Reader) {

	if strings.HasPrefix(bistroUrl, "file://") {
		bistroUrl := strings.Replace(bistroUrl, "file://", "", -1)
		documentReader = readFile(bistroUrl)
	} else {
		documentReader = readUrl(bistroUrl).Body
	}

	return documentReader
}

// retrieves a http response from the specivied url
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

// retrieves a file handle from the specified file path
func readFile(filePath string) *os.File {
	bistroPageReader, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	return bistroPageReader
}
