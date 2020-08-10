package restapi

import (
	"github.com/ansgarS/rate-my-bistro-crawler/config"
	"github.com/ansgarS/rate-my-bistro-crawler/jobs"
	"github.com/ansgarS/rate-my-bistro-crawler/persister"
	"github.com/yarf-framework/yarf"
	"io"
	"strings"
)

// Define a simple resource
type JobResource struct {
	yarf.Resource
}

// Define all routes for this resource
func addResourceEndpoints(y *yarf.Yarf) {
	y.Add("/jobs", new(JobResource))
	y.Add("/jobs/:jobId", new(JobResource))
}

// Implement the GET method
func (h *JobResource) Get(c *yarf.Context) error {
	jobId := c.Param("jobId")

	if jobId != "" {
		handleGetRequestJobIdParameter(c, jobId)
		return nil
	}

	c.RenderJSON(jobs.JobQueue)

	return nil
}

// Implement the POST method
func (h *JobResource) Post(c *yarf.Context) error {
	date := c.Request.Body

	if date != nil {
		handlePostRequestDateParameter(c, date)
		return nil
	}

	c.Status(400)

	return nil
}

func handleGetRequestJobIdParameter(c *yarf.Context, jobId string) {
	var job jobs.Job
	persister.ReadDocumentIfExists(config.Cfg.JobCollectionName, jobId, job)
	if job.Id == "" {
		c.Render("No job found for jobId " + jobId)
		c.Status(404)
	} else {
		c.RenderJSON(job)
	}
}

func handlePostRequestDateParameter(c *yarf.Context, dateReader io.ReadCloser) {
	buf := new(strings.Builder)
	io.Copy(buf, dateReader)
	date := buf.String()

	if date == "" {
		c.Status(400)
		c.Render("Note date playload found in request body")
		return
	}

	jobId := jobs.EnqueueJob(date)
	c.Status(201)
	c.Render(jobId)

}
