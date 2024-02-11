package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"

	"golang_backend_template/config"
	restful "golang_backend_template/internal/controller/restful"
	adapter "golang_backend_template/internal/infra/adapter"
	memo "golang_backend_template/internal/infra/memo"
	"golang_backend_template/internal/usecase/impl"
	"golang_backend_template/pkg/httpserver"
	"golang_backend_template/pkg/logger"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		l.Error(fmt.Errorf("app - Run - client.NewClientWithOpts: %w", err))
	}

	trainingJobManager := impl.NewTrainingJobManager(
		memo.NewTrainingJobsMemory(),
		adapter.NewDockerAdapter(cli),
		adapter.NewTwccAdapter(cfg.TWCC.APIKey),
	)

	inferenceJobManager := impl.NewInferenceJobManager(
		memo.NewInferenceJobsMemory(),
		adapter.NewTwccAdapter(cfg.TWCC.APIKey),
	)

	handler := gin.New()
	restful.SetupRouter(handler,
		l,
		trainingJobManager,
		inferenceJobManager)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
