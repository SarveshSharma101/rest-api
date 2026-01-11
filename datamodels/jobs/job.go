package datamodels

import "time"

type JobStatus string

const (
	JobPending  JobStatus = "PENDING"
	JobRunning  JobStatus = "RUNNING"
	JobComplete JobStatus = "COMPLETED"
)

type Job struct {
	ID        string    `json:"id"`
	Status    JobStatus `json:"status"`
	StartedAt time.Time `json:"started_at"`
}
