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
func addJobsResource(server *yarf.Yarf) {
	server.Add("/jobs", new(JobResource))
	server.Add("/jobs/:jobId", new(JobResource))
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
	c.Render("No body payload found, but a string date ('2001-12-31') as body is required")
	return nil
}

// Define the handler for a GET request with jobId parameter
func handleGetRequestJobIdParameter(c *yarf.Context, jobId string) {
	var job jobs.Job
	persister.ReadDocumentIfExists(config.Get().JobCollectionName, jobId, &job)
	if job.Id == "" {
		c.Status(404)
		c.Render("No job found for jobId " + jobId)
	} else {
		c.RenderJSON(job)
	}
}

// Define the handler for a POST request
func handlePostRequestDateParameter(c *yarf.Context, dateReader io.ReadCloser) {
	// Convert the request body to a string
	buf := new(strings.Builder)
	io.Copy(buf, dateReader)
	date := buf.String()

	if date == "" {
		c.Status(400)
		c.Render("No date playload found in request body")
		return
	}

	jobId := jobs.EnqueueJob(date)
	c.Status(201)
	c.Render(jobId)
}
