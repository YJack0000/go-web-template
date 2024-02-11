package memo

import (
	"fmt"
	"sync"

	"golang_backend_template/internal/usecase/entity"
)

type InferenceJobsMemory struct {
	mu   sync.Mutex
	jobs map[string]entity.InferenceJob
}

func NewInferenceJobsMemory() *InferenceJobsMemory {
	return &InferenceJobsMemory{
		jobs: make(map[string]entity.InferenceJob),
	}
}

func (r *InferenceJobsMemory) StoreInferenceJob(j entity.InferenceJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[j.Job.ID] = j
	return nil
}

func (r *InferenceJobsMemory) GetInferenceJob(id string) (entity.InferenceJob, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if j, ok := r.jobs[id]; ok {
		return j, nil
	}

	return entity.InferenceJob{}, fmt.Errorf("JobsMemory - GetInferenceJob - job not found")
}

func (r *InferenceJobsMemory) GetAllInferenceJob() ([]entity.InferenceJob, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	jobs := make([]entity.InferenceJob, 0, 64)

	for _, j := range r.jobs {
		jobs = append(jobs, j)
	}

	return jobs, nil
}

func (r *InferenceJobsMemory) DeleteInferenceJob(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.jobs, id)
	return nil
}
