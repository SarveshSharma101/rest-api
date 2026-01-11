package handler

import (
	"net/http"
	datamodels "rest-api/rest-api/datamodels/jobs"
	"rest-api/rest-api/utils"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var jobs map[string]datamodels.Job
var mu sync.RWMutex
var statusApi = "http://localhost:5000/job/{:jobid}"

func Job(ctx *gin.Context) {
	if len(jobs) == 0 {
		jobs = map[string]datamodels.Job{}
	}
	jobId, err := utils.GenerateUID()
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"err": err,
			},
		)
		return
	}

	job := datamodels.Job{
		ID:        jobId,
		Status:    datamodels.JobPending,
		StartedAt: time.Now(),
	}

	mu.Lock()
	jobs[jobId] = job
	mu.Unlock()

	go runJob(jobId)

	ctx.JSON(
		http.StatusCreated,
		gin.H{
			"Job-Id":    jobId,
			"Job":       job,
			"statusApi": statusApi,
		},
	)
}

func JobStatus(ctx *gin.Context) {
	jobId := ctx.Param("jobId")
	if jobId == "" {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"err": "Job id is required in path param",
			},
		)
	}

	mu.Lock()
	jobStatus := jobs[jobId].Status
	mu.Unlock()

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"Job Status": jobStatus,
		},
	)
}

func runJob(jobId string) {

	mu.Lock()
	job := jobs[jobId]
	job.Status = datamodels.JobRunning
	jobs[jobId] = job
	mu.Unlock()

	time.Sleep(10 * time.Second)

	mu.Lock()
	job.Status = datamodels.JobComplete
	jobs[jobId] = job
	mu.Unlock()

}
