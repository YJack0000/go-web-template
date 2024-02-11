package memo

import (
	"fmt"
	"sync"

	"golang_backend_template/internal/usecase/entity"
)

type TrainingJobsMemory struct {
	mu            sync.Mutex
	twccJobs      map[string]entity.TwccJob
	containerJobs map[string]entity.ContainerJob
	jobHistory    map[string]entity.GenericJob
}

func NewTrainingJobsMemory() *TrainingJobsMemory {
	return &TrainingJobsMemory{
		twccJobs:      make(map[string]entity.TwccJob),
		containerJobs: make(map[string]entity.ContainerJob),
		jobHistory:    make(map[string]entity.GenericJob),
	}
}

func (r *TrainingJobsMemory) PushTwccJob(j entity.TwccJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.twccJobs[j.Job.ID] = j
	return nil
}

func (r *TrainingJobsMemory) PushContainerJob(j entity.ContainerJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.containerJobs[j.Job.ID] = j
	return nil
}

func (r *TrainingJobsMemory) GetJob(id string) (entity.GenericJob, error) {
	if j, ok := r.twccJobs[id]; ok {
		return j.Job, nil
	}

	if j, ok := r.containerJobs[id]; ok {
		return j.Job, nil
	}

	if j, ok := r.jobHistory[id]; ok {
		return j, nil
	}

	return entity.GenericJob{}, fmt.Errorf("TrainingJobsMemory - GetJob - job not found")
}

func (r *TrainingJobsMemory) GetTwccJobList() ([]entity.TwccJob, error) {
	jobs := make([]entity.TwccJob, 0, 64)

	for _, j := range r.twccJobs {
		jobs = append(jobs, j)
	}

	return jobs, nil
}

func (r *TrainingJobsMemory) GetContainerJobList() ([]entity.ContainerJob, error) {
	jobs := make([]entity.ContainerJob, 0, 64)

	for _, j := range r.containerJobs {
		jobs = append(jobs, j)
	}

	return jobs, nil
}

func (r *TrainingJobsMemory) GetHistoryJobList() ([]entity.GenericJob, error) {
	jobs := make([]entity.GenericJob, 0, 64)

	for _, j := range r.jobHistory {
		jobs = append(jobs, j)
	}

	return jobs, nil
}

func (r *TrainingJobsMemory) storeHistory(job entity.GenericJob) {
	job.Status = "finished"
	r.jobHistory[job.ID] = job
}

func (r *TrainingJobsMemory) DeleteTwccJob(id string) error {
	if _, ok := r.twccJobs[id]; ok {
		r.storeHistory(r.twccJobs[id].Job)
		delete(r.twccJobs, id)
		return nil
	}

	return fmt.Errorf("TrainingJobsMemory - DeleteTwccJob - job not found")
}

func (r *TrainingJobsMemory) DeleteContainerJob(id string) error {
	if _, ok := r.containerJobs[id]; ok {
		r.storeHistory(r.containerJobs[id].Job)
		delete(r.containerJobs, id)
		return nil
	}

	return fmt.Errorf("TrainingJobsMemory - DeleteContainerJob - job not found")
}
