package adapter

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerAdapter struct {
	dockerClient *client.Client
}

func NewDockerAdapter(dockerClient *client.Client) *DockerAdapter {
	return &DockerAdapter{dockerClient: dockerClient}
}

func (r *DockerAdapter) pullImage(ctx context.Context, dockerImageName string) error {
	_, err := r.dockerClient.ImagePull(ctx, dockerImageName, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("DockerAdapter - PullImage - r.dockerClient.ImagePull: %w", err)
	}

	return nil
}

func (r *DockerAdapter) CreateContainer(ctx context.Context, dockerImageName string) (string, error) {
	if err := r.pullImage(ctx, dockerImageName); err != nil {
		return "", fmt.Errorf("DockerAdapter - CreateContainerJob - r.pullImage: %w", err)
	}
	resp, err := r.dockerClient.ContainerCreate(ctx, &container.Config{
		Image: dockerImageName,
	}, nil, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("DockerAdapter - CreateContainerJob - r.dockerClient.ContainerCreate: %w", err)
	}

	return resp.ID, nil
}

func (r *DockerAdapter) ContainerStartWithCallback(ctx context.Context, containerID string, callback func()) error {
	if err := r.dockerClient.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	statusCh, errCh := r.dockerClient.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("DockerAdapter - ContainerFinishedCallback - r.dockerClient.ContainerWait: %w", err)
		}
	case <-statusCh:
		callback()
	}

	return nil
}
