package ports

import (
	"context"
)

type (
	ContainerManager interface {
		CreateContainer(context.Context, string) (string, error)
		ContainerStartWithCallback(context.Context, string, func()) error
	}
)
