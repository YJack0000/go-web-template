package usecase

import "golang_backend_template/internal/usecase/entity"

type TrainingJobRequester interface {
	CreateJob(job entity.GenericJob, dockerImageName string, twccJobId string) error
	GetJob(id string) (entity.GenericJob, error)
	GetAllJobs() ([]entity.GenericJob, error)
	DeleteJob(jobID string) error
}
