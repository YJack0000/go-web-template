package ports

import "golang_backend_template/internal/usecase/entity"

type TrainingJobsRepo interface {
	PushContainerJob(entity.ContainerJob) error
	PushTwccJob(entity.TwccJob) error
	GetJob(string) (entity.GenericJob, error)
	GetTwccJobList() ([]entity.TwccJob, error)
	GetContainerJobList() ([]entity.ContainerJob, error)
	GetHistoryJobList() ([]entity.GenericJob, error)
	DeleteTwccJob(string) error
	DeleteContainerJob(string) error
}
