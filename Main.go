package main

import (
	"fmt"
	"github.com/ansgarS/rate-my-bistro-crawler/meals"
)

func main() {
	meals := meals.Start()
	fmt.Print(meals)
}
