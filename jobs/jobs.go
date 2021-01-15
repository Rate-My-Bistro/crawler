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
	Key          string   `json:"_key,omitempty"` // unique identifier for the database
	Id           string   `json:"id,omitempty"`   // uuid that unique identifies the job
	DateToParse  string   `json:"dateToParse"`    // The date which the parser should parse / has parsed.
	Status       string   `json:"status"`         // PENDING | RUNNING |  SUCCESS | FAILURE
	EnqueuedTime string   `json:"enqueuedTime"`   // the time the job was enqueued
	StartedTime  string   `json:"startedTime"`    // the time the job has started the parsing
	FinishedTime string   `json:"finishedTime"`   // the time the job has finished the parsing process
	Additional   []string `json:"additional"`     // optional information to keep near to the job (e.g. error messages)
}

// Holds all jobs in memory as a queue
var JobQueue = make([]Job, 0)

// Crates a new scheduler for the configured interval
func init() {
	schedulerTick := config.Get().JobSchedulerTickInSeconds

	s1 := gocron.NewScheduler(time.UTC)
	_, err := s1.Every(schedulerTick).Seconds().Do(processNextJob)
	s1.StartAsync()

	if err != nil {
		log.Fatal("Failed to init job scheduling", err)
	}
}

// Gets called on every tick of the scheduler
// It dequeues the head of the queue and start the parsing process.
// Every job status change is persisted to the job collection.
func processNextJob() {
	if len(JobQueue) <= 0 {
		return
	}

	// dequeue the next job and prepare it
	nextJob := DequeueJob()
	nextJob.StartedTime = time.Now().Format(time.RFC3339)
	nextJob.Status = "RUNNING"
	persister.PersistDocument(config.Get().JobCollectionName, nextJob)

	// start the meal crawling and store the result in the database
	log.Println("Start crawling meals for date " + nextJob.DateToParse)
	crawledMeals, err := crawler.CrawlAtDate(config.Get().BistroUrl, nextJob.DateToParse)
	if err != nil {
		jobFailureFinished(nextJob, err)
		return
	}
	persister.PersistDocuments(config.Get().MealCollectionName, ToIdentifiables(crawledMeals))

	// mark the job as finished successful
	jobSuccessFinished(nextJob)
	log.Println("Finished crawling meals for date " + nextJob.DateToParse)
}

func jobSuccessFinished(job Job) {
	job.FinishedTime = time.Now().Format(time.RFC3339)
	job.Status = "SUCCESS"
	persister.PersistDocument(config.Get().JobCollectionName, job)
}

func jobFailureFinished(job Job, err error) {
	job.FinishedTime = time.Now().Format(time.RFC3339)
	job.Status = "FAILURE"
	job.Additional = []string{err.Error()}
	persister.PersistDocument(config.Get().JobCollectionName, job)
}

// Enqueues a new parser job for a specific date at the end of the queue
// Returns the id of the created job
func EnqueueJob(dateToParse string) string {
	uid, _ := uuid.NewV4()
	identifier := uid.String()
	newJob := Job{
		Key:          identifier,
		Id:           identifier,
		Status:       "PENDING",
		EnqueuedTime: time.Now().Format(time.RFC3339),
		DateToParse:  dateToParse,
	}
	JobQueue = append(JobQueue, newJob)
	persister.PersistDocument(config.Get().JobCollectionName, newJob)
	return identifier
}

// Dequeues the head of the queue.
// This removes the dequeued item
func DequeueJob() (nextJob Job) {
	nextJob = JobQueue[0]
	JobQueue = JobQueue[1:] // Discard top element
	return nextJob
}

// Identifiable interface implantation for the struct job
func (job Job) GetId() string {
	return job.Key
}

// Casts a slice of meals to the a slice of Identifiable interfaces
func ToIdentifiables(meals []crawler.Meal) []persister.Identifiable {
	identifiables := make([]persister.Identifiable, len(meals))
	for i := range meals {
		identifiables[i] = meals[i]
	}
	return identifiables
}

// Removes all entries from the job queue
func RemoveAllJobs() {
	JobQueue = make([]Job, 0)
}
