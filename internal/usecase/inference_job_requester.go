package usecase

import "golang_backend_template/internal/usecase/entity"

type InferenceJobRequester interface {
	CreateJob(job entity.GenericJob) (string, error)
	GetJob(id string) (entity.GenericJob, error)
	GetAllJobs() ([]entity.GenericJob, error)
	DeleteJob(jobID string) error
}
