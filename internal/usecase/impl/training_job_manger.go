package impl

import (
	"context"
	"fmt"
	"time"

	"golang_backend_template/internal/usecase/entity"
	"golang_backend_template/internal/usecase/ports"
)

type TrainingJobManager struct {
	repo   ports.TrainingJobsRepo
	docker ports.ContainerManager
	twcc   ports.TwccManager
}

func NewTrainingJobManager(m ports.TrainingJobsRepo, d ports.ContainerManager, w ports.TwccManager) *TrainingJobManager {
	return &TrainingJobManager{
		repo:   m,
		docker: d,
		twcc:   w,
	}
}

func (uc *TrainingJobManager) CreateJob(job entity.GenericJob, dockerImageName string, twccJobId string) error {
	containerJobs, err := uc.repo.GetContainerJobList()

	if err != nil {
		return fmt.Errorf("TrainingJobManager - CreateJob - s.repo.GetTwccJobList: %w", err)
	}

	if len(containerJobs) < 2 {
		job.Status = "running on docker"
		err = uc.repo.PushContainerJob(entity.ContainerJob{
			Job:             job,
			DockerImageName: dockerImageName,
		})

		if err != nil {
			return fmt.Errorf("TrainingJobManager - CreateJob - s.repo.CreateContainerJob: %w", err)
		}

		ctx := context.Background()
		containerID, err := uc.docker.CreateContainer(ctx, dockerImageName)
		if err != nil {
			uc.repo.DeleteContainerJob(job.ID)
			return fmt.Errorf("TrainingJobManager - CreateJob - s.docker.CreateContainerJob: %w", err)
		}

		// function to remove container job from queue after container job is done
		go uc.docker.ContainerStartWithCallback(ctx, containerID, func() {
			err := uc.repo.DeleteContainerJob(job.ID)
			if err != nil {
				fmt.Errorf("TrainingJobManager - CreateJob - s.repo.DeleteContainerJob: %w", err)
			}
		})

		return nil
	}

	job.Status = "running on twcc"
	err = uc.repo.PushTwccJob(entity.TwccJob{
		Job:       job,
		TwccJobId: twccJobId,
	})
	if err != nil {
		return fmt.Errorf("TrainingJobManager - CreateJob - s.repo.CreateTwccJob: %w", err)
	}

	err = uc.twcc.RunTwccJob(twccJobId)
	if err != nil {
		uc.repo.DeleteTwccJob(job.ID)
		return fmt.Errorf("TrainingJobManager - CreateJob - s.twcc.RunTwccJob: %w", err)
	}

	// TODO: more elegant way to check if twcc job is done
	go func() {
		count := 0
		for ; count < 10; count++ {
			time.Sleep(3 * time.Second)
			status, _ := uc.twcc.GetTwccJobStatus(twccJobId)
			// if err != nil {
			// 	fmt.Errorf("TrainingJobManager - CreateJob - s.twcc.GetTwccJobStatus: %w", err)
			// }

			if status == "Inactive" {
				_ = uc.repo.DeleteTwccJob(job.ID)
				// if err != nil {
				// 	fmt.Errorf("TrainingJobManager - CreateJob - s.repo.DeleteTwccJob: %w", err)
				// }
				break
			}
		}
	}()

	return nil
}

func (uc *TrainingJobManager) GetJob(id string) (entity.GenericJob, error) {
	return uc.repo.GetJob(id)
}

func (uc *TrainingJobManager) GetAllJobs() ([]entity.GenericJob, error) {
	containerJobs, err := uc.repo.GetContainerJobList()
	if err != nil {
		return nil, fmt.Errorf("TrainingJobManager - GetAllJob - s.repo.GetTwccJobList: %w", err)
	}

	twccJobs, err := uc.repo.GetTwccJobList()
	if err != nil {
		return nil, fmt.Errorf("TrainingJobManager - GetAllJob - s.repo.GetTwccJobList: %w", err)
	}

	historyJobs, err := uc.repo.GetHistoryJobList()
	if err != nil {
		return nil, fmt.Errorf("TrainingJobManager - GetAllJob - s.repo.GetHistoryJobList: %w", err)
	}

	jobs := make([]entity.GenericJob, 0, 64)

	for _, j := range containerJobs {
		jobs = append(jobs, j.Job)
	}

	for _, j := range twccJobs {
		jobs = append(jobs, j.Job)
	}

	for _, j := range historyJobs {
		jobs = append(jobs, j)
	}

	return jobs, nil
}

func (uc *TrainingJobManager) DeleteJob(id string) error {
	containerJobs, err := uc.repo.GetContainerJobList()

	if err != nil {
		return fmt.Errorf("TrainingJobManager - DeleteJob - s.repo.GetTwccJobList: %w", err)
	}

	if len(containerJobs) > 0 {
		err = uc.repo.DeleteContainerJob(id)
		if err != nil {
			return fmt.Errorf("TrainingJobManager - DeleteJob - s.repo.DeleteContainerJob: %w", err)
		}

		return nil
	}

	err = uc.repo.DeleteTwccJob(id)

	if err != nil {
		return fmt.Errorf("TrainingJobManager - DeleteJob - s.repo.DeleteTwccJob: %w", err)
	}

	return nil
}
