package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "golang_backend_template/docs"
	v1 "golang_backend_template/internal/controller/restful/v1"
	"golang_backend_template/internal/usecase"
	"golang_backend_template/pkg/logger"
)

// @title swagger test
// @version 1.0
// @description swagger test example
// @schemes http https
// @BasePath /v1
func SetupRouter(handler *gin.Engine, l logger.Interface, trainingJobManager usecase.TrainingJobRequester, inferenceJobManager usecase.InferenceJobRequester) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	handler.GET("/swagger/*any", swaggerHandler)
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.Group("/v1")
	{
		v1.InitTrainingJobRoutes(h, trainingJobManager, l)
		v1.InitInferenceJobRoutes(h, inferenceJobManager, l)
	}
}
