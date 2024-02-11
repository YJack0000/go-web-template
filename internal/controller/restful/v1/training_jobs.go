package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"golang_backend_template/internal/usecase"
	"golang_backend_template/internal/usecase/entity"
	"golang_backend_template/pkg/logger"
)

type TrainingJobController struct {
	u usecase.TrainingJobRequester
	l logger.Interface
}

func InitTrainingJobRoutes(handler *gin.RouterGroup, u usecase.TrainingJobRequester, l logger.Interface) {
	r := &TrainingJobController{u, l}

	h := handler.Group("/training-jobs")
	{
		h.GET("/all", r.list)
		h.POST("/create", r.create)
		h.GET(":id", r.get)
	}
}

type createTrainingJobRequest struct {
	TwccJobId       string `json:"twccJobId" example:"237139"`
	DockerImageName string `json:"dockerImageName" example:"yjack0000cs12/llm-training:latest"`
}

type createTrainingJobResponse struct {
	JobId string `json:"jobId" example:"12345"`
}

// @Summary     create training job
// @Description create training job
// @Tags  	    training-jobs
// @Accept      json
// @Produce     json
// @Success     200 {object} createTrainingJobResponse
// @Failure     500 {object} eResponse
// @Router      /training-jobs/create [post]
// @Param       req body createTrainingJobRequest true "request"
func (r *TrainingJobController) create(c *gin.Context) {
	var req createTrainingJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - create")
		errorResponse(c, 400, "invalid request")

		return
	}

	job := entity.GenericJob{
		ID:     uuid.New().String(),
		Name:   req.DockerImageName + "-" + req.TwccJobId,
		Status: "created",
	}

	err := r.u.CreateJob(job, req.DockerImageName, req.TwccJobId)
	if err != nil {
		r.l.Error(err, "http - v1 - create")
		errorResponse(c, 500, "database problems")

		return
	}

	c.JSON(200, createTrainingJobResponse{job.ID})
}

type getJobResponse struct {
	Job entity.GenericJob `json:"job"`
}

// @Summary     get training job
// @Description get training job
// @Tags  	    training-jobs
// @Accept      json
// @Produce     json
// @Param       id   path      string  true  "Job ID"
// @Success     200 {object} getJobResponse
// @Failure     500 {object} eResponse
// @Router      /training-jobs/{id} [get]
func (r *TrainingJobController) get(c *gin.Context) {
	id := c.Param("id")

	job, err := r.u.GetJob(id)
	if err != nil {
		r.l.Error(err, "http - v1 - get")
		errorResponse(c, 500, "database problems")

		return
	}

	c.JSON(200, getJobResponse{job})
}

type listTrainingJobResponse struct {
	Jobs []entity.GenericJob `json:"jobs"`
}

// @Summary     list all training jobs
// @Description list all training jobs
// @Tags  	    training-jobs
// @Accept      json
// @Produce     json
// @Success     200 {object} getJobResponse
// @Failure     500 {object} eResponse
// @Router      /training-jobs/all [get]
func (r *TrainingJobController) list(c *gin.Context) {
	jobs, err := r.u.GetAllJobs()
	if err != nil {
		r.l.Error(err, "http - v1 - get")
		errorResponse(c, 500, "database problems")

		return
	}

	c.JSON(200, listTrainingJobResponse{jobs})
}
