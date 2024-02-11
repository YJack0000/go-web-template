package usecase

import (
	"context"
	"fmt"
	"time"

	"golang_backend_template/internal/usecase/entity"
	"golang_backend_template/internal/usecase/ports"
)

type TrainingJobManager struct {
	memo   ports.TrainingJobsRepo
	docker ports.ContainerManager
	twcc   ports.TwccManager
}

func NewTrainingJobManager(m ports.TrainingJobsRepo, d ports.ContainerManager, w ports.TwccManager) *TrainingJobManager {
	return &TrainingJobManager{
		memo:   m,
		docker: d,
		twcc:   w,
	}
}

func (uc *TrainingJobManager) CreateJob(job entity.GenericJob, dockerImageName string, twccJobId string) error {
	containerJobs, err := uc.memo.GetContainerJobList()

	if err != nil {
		return fmt.Errorf("TrainingJobManager - CreateJob - s.memo.GetTwccJobList: %w", err)
	}

	if len(containerJobs) < 2 {
		job.Status = "running on docker"
		err = uc.memo.PushContainerJob(entity.ContainerJob{
			Job:             job,
			DockerImageName: dockerImageName,
		})

		if err != nil {
			return fmt.Errorf("TrainingJobManager - CreateJob - s.memo.CreateContainerJob: %w", err)
		}

		ctx := context.Background()
		containerID, err := uc.docker.CreateContainer(ctx, dockerImageName)
		if err != nil {
			uc.memo.DeleteContainerJob(job.ID)
			return fmt.Errorf("TrainingJobManager - CreateJob - s.docker.CreateContainerJob: %w", err)
		}

		// function to remove container job from queue after container job is done
		go uc.docker.ContainerStartWithCallback(ctx, containerID, func() {
			err := uc.memo.DeleteContainerJob(job.ID)
			if err != nil {
				fmt.Errorf("TrainingJobManager - CreateJob - s.memo.DeleteContainerJob: %w", err)
			}
		})

		return nil
	}

	job.Status = "running on twcc"
	err = uc.memo.PushTwccJob(entity.TwccJob{
		Job:       job,
		TwccJobId: twccJobId,
	})
	if err != nil {
		return fmt.Errorf("TrainingJobManager - CreateJob - s.memo.CreateTwccJob: %w", err)
	}

	err = uc.twcc.RunTwccJob(twccJobId)
	if err != nil {
		uc.memo.DeleteTwccJob(job.ID)
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
				_ = uc.memo.DeleteTwccJob(job.ID)
				// if err != nil {
				// 	fmt.Errorf("TrainingJobManager - CreateJob - s.memo.DeleteTwccJob: %w", err)
				// }
				break
			}
		}
	}()

	return nil
}

func (uc *TrainingJobManager) GetJob(id string) (entity.GenericJob, error) {
	return uc.memo.GetJob(id)
}

func (uc *TrainingJobManager) GetAllJob() ([]entity.GenericJob, error) {
	containerJobs, err := uc.memo.GetContainerJobList()
	if err != nil {
		return nil, fmt.Errorf("TrainingJobManager - GetAllJob - s.memo.GetTwccJobList: %w", err)
	}

	twccJobs, err := uc.memo.GetTwccJobList()
	if err != nil {
		return nil, fmt.Errorf("TrainingJobManager - GetAllJob - s.memo.GetTwccJobList: %w", err)
	}

	historyJobs, err := uc.memo.GetHistoryJobList()
	if err != nil {
		return nil, fmt.Errorf("TrainingJobManager - GetAllJob - s.memo.GetHistoryJobList: %w", err)
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
	containerJobs, err := uc.memo.GetContainerJobList()

	if err != nil {
		return fmt.Errorf("TrainingJobManager - DeleteJob - s.memo.GetTwccJobList: %w", err)
	}

	if len(containerJobs) > 0 {
		err = uc.memo.DeleteContainerJob(id)
		if err != nil {
			return fmt.Errorf("TrainingJobManager - DeleteJob - s.memo.DeleteContainerJob: %w", err)
		}

		return nil
	}

	err = uc.memo.DeleteTwccJob(id)

	if err != nil {
		return fmt.Errorf("TrainingJobManager - DeleteJob - s.memo.DeleteTwccJob: %w", err)
	}

	return nil
}
