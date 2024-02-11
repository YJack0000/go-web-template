package impl

import (
	"fmt"
	"time"

	"golang_backend_template/internal/usecase/entity"
	"golang_backend_template/internal/usecase/ports"
)

type InferenceJobManager struct {
	repo ports.InferenceJobRepo
	twcc ports.TwccManager
}

func NewInferenceJobManager(m ports.InferenceJobRepo, t ports.TwccManager) *InferenceJobManager {
	return &InferenceJobManager{
		repo: m,
		twcc: t,
	}
}

func (uc *InferenceJobManager) CreateJob(job entity.GenericJob) (string, error) {
	twccCCSId, err := uc.twcc.CreateTwccCCS()
	if err != nil {
		return "", fmt.Errorf("InferenceJobManager - CreateJob - s.twcc.CreateTwccCCS: %w", err)
	}

	time.Sleep(1 * time.Second / 2)
	err = uc.twcc.TwccCCSAssociateIP(twccCCSId)
	if err != nil {
		return "", fmt.Errorf("InferenceJobManager - CreateJob - s.twcc.TwccCCSAssociateIP: %w", err)
	}

	entryPoint, err := uc.twcc.GetTwccCCSEntryPoint(twccCCSId)
	if err != nil {
		return "", fmt.Errorf("InferenceJobManager - CreateJob - s.twcc.GetTwccCCSEntryPoint: %w", err)
	}

	job.Status = "inference running on twcc"
	err = uc.repo.StoreInferenceJob(entity.InferenceJob{Job: job, TwccCCSId: twccCCSId, EntryPoint: entryPoint})
	if err != nil {
		return "", fmt.Errorf("InferenceJobManager - CreateJob - s.repo.CreateInferenceJob: %w", err)
	}

	return entryPoint, nil
}

func (uc *InferenceJobManager) GetJob(id string) (entity.GenericJob, error) {
	job, err := uc.repo.GetInferenceJob(id)
	if err != nil {
		return entity.GenericJob{}, fmt.Errorf("InferenceJobManager - GetJob - s.repo.GetInferenceJob: %w", err)
	}

	return job.Job, nil
}

func (uc *InferenceJobManager) GetAllJobs() ([]entity.GenericJob, error) {
	jobs, err := uc.repo.GetAllInferenceJob()
	if err != nil {
		return nil, fmt.Errorf("InferenceJobManager - GetAllJob - s.repo.GetAllInferenceJob: %w", err)
	}

	genericJobs := make([]entity.GenericJob, 0, 64)
	for _, j := range jobs {
		genericJobs = append(genericJobs, j.Job)
	}

	return genericJobs, nil
}

func (uc *InferenceJobManager) getTwccCCSId(id string) (string, error) {
	job, err := uc.repo.GetInferenceJob(id)
	if err != nil {
		return "", fmt.Errorf("InferenceJobManager - GetTwccJobId - s.repo.GetInferenceJob: %w", err)
	}

	return job.TwccCCSId, nil
}

func (uc *InferenceJobManager) DeleteJob(id string) error {
	twccCCSId, err := uc.getTwccCCSId(id)

	if err != nil {
		return fmt.Errorf("InferenceJobManager - DeleteJob - s.getTwccCCSId: %w", err)
	}

	err = uc.twcc.DeleteTwccCCS(twccCCSId)
	if err != nil {
		return fmt.Errorf("InferenceJobManager - DeleteJob - s.twcc.DeleteTwccCCS: %w", err)
	}

	err = uc.repo.DeleteInferenceJob(id)
	if err != nil {
		return fmt.Errorf("InferenceJobManager - DeleteJob - s.repo.DeleteInferenceJob: %w", err)
	}

	return nil
}
