package restapi

import (
	"github.com/Rate-My-Bistro/crawler/config"
	"github.com/Rate-My-Bistro/crawler/jobs"
	"github.com/Rate-My-Bistro/crawler/persister"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
	"time"
)

// See Declarative Comments Format: https://swaggo.github.io/swaggo.io/declarative_comments_format/general_api_info.html

// jobGet godoc
// @Summary Get all running job
// @Description get job all running jobs
// @Tags jobs
// @Accept plain/text
// @Success 200 {array} jobs.Job
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /jobs [get]
func jobGet() func(context *gin.Context) {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, jobs.JobQueue)
	}
}

// jobGet godoc
// @Summary Retrieve a job by it's id
// @Description get job by ID
// @Tags jobs
// @Accept plain/text
// @Produce application/json
// @Param id path string true "Job ID"
// @Success 200 {object} jobs.Job
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /jobs/{id} [get]
func jobGetWithParameter() func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id != "" {
			handleGetWithJobIdParameter(c, id)
		}
	}
}

// jobGet godoc
// @Summary Create a new parser job
// @Description create a new parser job for the specified date
// @Tags jobs
// @Produce plain/text
// @Accept plain/text
// @Param date body string true "Date to parse in yyyy-mm-dd" format date
// @Success 201 {string} string
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /jobs [post]
func jobPost() func(context *gin.Context) {
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
