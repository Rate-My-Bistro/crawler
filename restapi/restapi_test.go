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

//TODO create a test where the state of an meal job goes to error
func TestPostJobWithDateInNextFuture(t *testing.T) {
	router := setupRouter()

	// When posting a new job with an invalid date
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/jobs", toReader("13-08-2020"))
	router.ServeHTTP(resp, req)
	jobId := resp.Body.String()

	// In first place it should be fine
	assert.Equal(t, 201, resp.Code)

	// But when we wait until the job was started
	time.Sleep(1 * time.Second)

	// it should have the status failed
	resp = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/jobs/"+jobId, nil)
	router.ServeHTTP(resp, req)

	s := toJson(t, resp.Body.String())
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, s["id"], jobId)
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
