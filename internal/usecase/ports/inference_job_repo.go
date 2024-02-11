package ports

import "golang_backend_template/internal/usecase/entity"

type InferenceJobRepo interface {
	StoreInferenceJob(entity.InferenceJob) error
	GetInferenceJob(string) (entity.InferenceJob, error)
	GetAllInferenceJob() ([]entity.InferenceJob, error)
	DeleteInferenceJob(string) error
}
