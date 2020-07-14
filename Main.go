package main

import (
	"fmt"
	"github.com/ansgarS/rate-my-bistro-crawler/crawler"
	"log"
	"net/http"
)

const bistroUrl = "https://bistro.cgm.ag/index.php"

func main() {
	httpResponse := readUrl()

	meals := crawler.Start(httpResponse.Body)

	fmt.Print(meals)
}

func readUrl() *http.Response {
	res, err := http.Get(bistroUrl)

	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	return res
}
