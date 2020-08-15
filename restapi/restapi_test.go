package restapi

import (
	"bytes"
	"encoding/json"
	"github.com/ansgarS/rate-my-bistro-crawler/jobs"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAllEndpointsInPositiveCase(t *testing.T) {
	router := setupRouter()

	// POST a jew job
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/jobs", toReader("2020-08-13"))
	router.ServeHTTP(resp, req)
	jobId := resp.Body.String()

	assert.Equal(t, 201, resp.Code)
	assert.NotEmpty(t, jobId)

	// GET all running jobs
	resp = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/jobs", nil)
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, toJsonString(t, jobs.JobQueue), resp.Body.String())

	// GET this job by its id
	resp = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/jobs/"+jobId, nil)
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, resp.Body.String(), jobId)

	// Cleanup
	jobs.RemoveAllJobs()
}

func TestUnknownId(t *testing.T) {
	router := setupRouter()

	// GET this job by its id
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/jobs", nil)
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "[]", resp.Body.String())
}

func TestEmptyQueue(t *testing.T) {
	router := setupRouter()

	// GET all running jobs
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/jobs", nil)
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "[]", resp.Body.String())
}

func TestPostJobWithInvalidDate(t *testing.T) {
	router := setupRouter()

	// When posting a new job with an invalid date
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/jobs", toReader("13-08-2020"))
	router.ServeHTTP(resp, req)

	// Then the response should point out the mistake
	assert.Equal(t, 400, resp.Code)
}

func TestPostJobWithDateInPast(t *testing.T) {
	router := setupRouter()

	// When posting a new job with an invalid date
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/jobs", toReader("1900-01-01"))
	router.ServeHTTP(resp, req)
	jobId := resp.Body.String()

	// The job should be enqueued
	assert.Equal(t, 201, resp.Code)

	// lets wait until the job was processed
	time.Sleep(1 * time.Second)

	// retrieve the job
	resp = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/jobs/"+jobId, nil)
	router.ServeHTTP(resp, req)

	// then the job should be failed
	s := toJson(t, resp.Body.String())
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, s["id"], jobId)
	assert.Equal(t, s["status"], "FAILURE")
}

func toReader(s string) io.Reader {
	return bytes.NewBufferString(s)
}

func toJson(t *testing.T, jsonString string) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	return result
}

func toJsonString(t *testing.T, input interface{}) string {
	b, err := json.Marshal(input)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	return string(b)
}
