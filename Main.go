package main

import (
	"fmt"
	"github.com/ansgarS/rate-my-bistro-crawler/jobs"
	"time"
)

// the application cycle
func main() {
	jobs.EnqueueJob("2020-05-04")
	jobs.EnqueueJob("2020-04-04")
	jobs.EnqueueJob("2020-03-04")

	go forever()
	//Keep this goroutine from exiting
	//so that the program doesn't end.
	_, _ = fmt.Scanln()
}

func forever() {
	for {
		//goroutine that run indefinitely.
		time.Sleep(time.Second)
	}
}
