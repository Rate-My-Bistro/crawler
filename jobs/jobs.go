package jobs

/*
Package jobs implements a service that is able to enqueue crawler jobs.
This jobs gets dequeued periodically and their status is persisted into the jobs document collection.
*/
import (
	"github.com/ansgarS/rate-my-bistro-crawler/config"
	"github.com/ansgarS/rate-my-bistro-crawler/crawler"
	"github.com/ansgarS/rate-my-bistro-crawler/persister"
	"github.com/go-co-op/gocron"
	"github.com/nu7hatch/gouuid"
	"log"
	"time"
)

// Represents a crawler job
type Job struct {
	Id          string   `json:"_key,omitempty"`
	DateToParse string   `json:"dateToParse"`
	Status      string   `json:"status"` // PENDING | RUNNING |  SUCCESS | FAILURE
	Enqueued    string   `json:"enqueued"`
	Started     string   `json:"started"`
	Finished    string   `json:"finished"`
	Additional  []string `json:"additional"`
}

// Holds all jobs in memory as a queue
var jobQueue = make([]Job, 0)

// Crates a new scheduler for the configured interval
func init() {
	s1 := gocron.NewScheduler(time.UTC)
	_, err := s1.Every(config.Cfg.JobSchedulerTickInSeconds).Seconds().Do(processNextJob)
	s1.StartAsync()

	if err != nil {
		log.Fatal("Failed to init job scheduling", err)
	}
}

// Gets called on every tick of the scheduler
// It dequeues the head of the queue and start the parsing process.
// Every job status change is persisted to the job collection.
func processNextJob() {
	if len(jobQueue) <= 0 {
		return
	}

	// dequeue the next job and prepare it
	nextJob := DequeueJob()
	nextJob.Started = time.Now().Format(time.RFC3339)
	nextJob.Status = "RUNNING"
	persister.PersistDocument(config.Cfg.JobCollectionName, nextJob)

	// start the meal crawling and store the result in the database
	log.Println("Start crawling meals for date " + nextJob.DateToParse)
	crawledMeals := crawler.CrawlAtDate(config.Cfg.BistroUrl, nextJob.DateToParse)
	persister.PersistDocuments(config.Cfg.MealCollectionName, ToIdentifiables(crawledMeals))

	// mark the job as finished successful
	nextJob.Finished = time.Now().Format(time.RFC3339)
	nextJob.Status = "SUCCESS"
	persister.PersistDocument(config.Cfg.JobCollectionName, nextJob)
	log.Println("Finished crawling meals for date " + nextJob.DateToParse)
}

// Enqueues a new parser job for a specific date at the end of the queue
func EnqueueJob(dateToParse string) {
	uid, _ := uuid.NewV4()
	newJob := Job{
		Id:          uid.String(),
		Status:      "PENDING",
		Enqueued:    time.Now().Format(time.RFC3339),
		DateToParse: dateToParse,
	}
	jobQueue = append(jobQueue, newJob)
	persister.PersistDocument(config.Cfg.JobCollectionName, newJob)
}

// Dequeues the head of the queue.
// This removes the dequeued item
func DequeueJob() (nextJob Job) {
	nextJob = jobQueue[0]
	jobQueue = jobQueue[1:] // Discard top element
	return nextJob
}

// Identifiable interface implantation for the struct job
func (job Job) GetId() string {
	return job.Id
}

// Casts a slice of meals to the a slice of Identifiable interfaces
func ToIdentifiables(meals []crawler.Meal) []persister.Identifiable {
	identifiables := make([]persister.Identifiable, len(meals))
	for i := range meals {
		identifiables[i] = meals[i]
	}
	return identifiables
}
