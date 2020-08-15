package restapi

import (
	"github.com/ansgarS/rate-my-bistro-crawler/config"
	"github.com/ansgarS/rate-my-bistro-crawler/jobs"
	"github.com/ansgarS/rate-my-bistro-crawler/persister"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
	"time"
)

// Define all routes for this resource
func addJobsResource(gin *gin.Engine) {
	gin.GET("/jobs", Get())
	gin.GET("/jobs/:jobId", GetWithParameter())

	gin.POST("/jobs", Post())
}

func Get() func(context *gin.Context) {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, jobs.JobQueue)
	}
}

func GetWithParameter() func(c *gin.Context) {
	return func(c *gin.Context) {
		jobId := c.Param("jobId")
		if jobId != "" {
			handleGetWithJobIdParameter(c, jobId)
		}
	}
}

func Post() func(context *gin.Context) {
	return func(context *gin.Context) {
		date := context.Request.Body

		if date == nil {
			context.String(http.StatusBadRequest, "No body payload found, but a string date ('2001-12-31') as body is required")
			return
		}

		handlePostWithBodyParam(context, date)
	}
}

// Define the handler for a GET request with jobId parameter
func handleGetWithJobIdParameter(c *gin.Context, jobId string) {
	var job jobs.Job
	persister.ReadDocumentIfExists(config.Get().JobCollectionName, jobId, &job)
	if job.Id == "" {
		c.String(404, "No job found for jobId "+jobId)
	} else {
		c.JSON(http.StatusOK, job)
	}
}

// Define the handler for a POST request
func handlePostWithBodyParam(c *gin.Context, dateReader io.ReadCloser) {
	// Convert the request body to a string
	buf := new(strings.Builder)
	io.Copy(buf, dateReader)
	date := buf.String()

	if date == "" {
		c.String(http.StatusBadRequest, "No date payload found in request body")
		return
	}

	_, err := time.Parse("2006-01-02", date)
	if err == nil {
		jobId := jobs.EnqueueJob(date)
		c.String(http.StatusCreated, jobId)
	} else {
		c.String(http.StatusBadRequest, "Invalid date format, expected was 'yyyy-mm-dd' but got "+date)
	}
}
