package jobs

import (
	"github.com/ansgarS/rate-my-bistro-crawler/config"
	"github.com/ansgarS/rate-my-bistro-crawler/crawler"
	"github.com/ansgarS/rate-my-bistro-crawler/persister"
	"github.com/jasonlvhit/gocron"
	"github.com/nu7hatch/gouuid"
	"time"
)

type Job struct {
	Id          string   `json:"_key,omitempty"`
	DateToParse string   `json:"dateToParse"`
	Status      string   `json:"status"` // PENDING | RUNNING |  SUCCESS | FAILURE
	Enqueued    string   `json:"enqueued"`
	Started     string   `json:"started"`
	Finished    string   `json:"finished"`
	Additional  []string `json:"additional"`
}

var jobQueue = make([]Job, 0)

func init() {
	gocron.Every(10).Second().Do(processNextJob)
}

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
	crawledMeals := crawler.CrawlAtDate(config.Cfg.BistroUrl, nextJob.DateToParse)
	persister.PersistDocuments(config.Cfg.MealCollectionName, ToIdentifiables(crawledMeals))

	// mark the job as finished successful
	nextJob.Finished = time.Now().Format(time.RFC3339)
	nextJob.Status = "SUCCESS"
	persister.PersistDocument(config.Cfg.JobCollectionName, nextJob)
}

// Enqueues a new parser job for a specific date
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

func DequeueJob() (nextJob Job) {
	nextJob = jobQueue[0]
	jobQueue = jobQueue[1:] // Discard top element
	return nextJob
}

func ClearJobs() {
	jobQueue = make([]Job, 0)
}

func (job Job) GetId() string {
	return job.Id
}

func ToIdentifiables(meals []crawler.Meal) []persister.Identifiable {
	identifiables := make([]persister.Identifiable, len(meals))
	for i := range meals {
		identifiables[i] = meals[i]
	}
	return identifiables
}
