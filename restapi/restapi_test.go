package restapi

import (
	"encoding/json"
	"github.com/ansgarS/rate-my-bistro-crawler/config"
	"github.com/ansgarS/rate-my-bistro-crawler/jobs"
	"github.com/go-resty/resty/v2"
	"testing"
)

var JobsEndpoint = "http://" + config.Get().RestApiAddress + "/jobs"

func TestIssuingAJobThroughTheRestApi(t *testing.T) {
	//start http server async
	go func() {
		Serve()
	}()

	//setup client
	client := resty.New()

	t.Run("when multiple jobs was started we can receive all", func(t *testing.T) {
		_, err1 := startJob(nil, client, "2020-08-11")
		_, err2 := startJob(nil, client, "2020-08-12")
		_, err3 := startJob(nil, client, "2020-08-13")

		if err1 != nil {
			t.Error(err1)
		}
		if err2 != nil {
			t.Error(err2)
		}
		if err3 != nil {
			t.Error(err3)
		}

		retrievedJobs, err := getJobs(client)

		if err != nil {
			t.Error(err)
		}
		if len(retrievedJobs) != 3 {
			t.Errorf("wanted 3 jobs but got %d", len(retrievedJobs))
		}
	})

	t.Run("when a job was started we can receive the job by the id", func(t *testing.T) {
		jobId, err := startJob(t, client, "2020-08-11")

		if err != nil {
			t.Error(err)
		}

		job, err := getJob(client, jobId)

		if err != nil {
			t.Error(err)
		}
		if jobId != job.GetId() {
			t.Errorf("wanted job id %s but got %s", jobId, job.GetId())
		}
	})

	t.Run("when requesting a non existent job an error is returned", func(t *testing.T) {
		respCode, err := getJobStatusCode(client, "unknown-job")

		if err != nil {
			t.Error(err)
		}
		if respCode != 404 {
			t.Errorf("Expected status code 404 but got %d", respCode)
		}
	})
}

// creates a new job and returns it's jobId as string
func startJob(t *testing.T, client *resty.Client, date string) (string, error) {
	resp, err := client.R().
		SetBody(date).
		Post(JobsEndpoint)

	if resp.StatusCode() != 201 {
		t.Errorf("Expected status code 201 but got %d", resp.StatusCode())
	}

	return resp.String(), err
}

// gets a job by it's id and returns a json string
func getJob(client *resty.Client, jobId string) (jobs.Job, error) {
	resp, err := client.R().
		Get(JobsEndpoint + "/" + jobId)
	job := jobs.Job{}
	err = json.Unmarshal(resp.Body(), &job)
	return job, err
}

// gets a job by it's id and returns a json string
func getJobStatusCode(client *resty.Client, jobId string) (int, error) {
	resp, err := client.R().
		Get(JobsEndpoint + "/" + jobId)
	return resp.StatusCode(), err
}

// gets a job by it's id and returns a json string
func getJobs(client *resty.Client) ([]jobs.Job, error) {
	resp, err := client.R().
		Get(JobsEndpoint)
	var jobs []jobs.Job
	err = json.Unmarshal(resp.Body(), &jobs)
	return jobs, err
}
