package v1

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"golang_backend_template/internal/usecase"
	"golang_backend_template/internal/usecase/entity"
	"golang_backend_template/pkg/logger"
)

type InferenceJobController struct {
	u usecase.InferenceJobManager
	l logger.Interface
}

func InitInferenceJobRoutes(handler *gin.RouterGroup, u usecase.InferenceJobManager, l logger.Interface) {
	c := &InferenceJobController{u, l}

	h := handler.Group("/inference-jobs")
	{
		h.GET("", c.list)
		h.POST("", c.create)
		h.GET(":id", c.get)
		h.DELETE(":id", c.delete)
	}
}

type createInferenceJobRequest struct {
}

type createInferenceJobResponse struct {
	JobId      string `json:"jobId" example:"12345"`
	EntryPoint string `json:"entryPoint" example:"http://localhost:8080"`
}

// @Summary     create inference job
// @Description create inference job
// @ID          create
// @Tags  	    inference-jobs
// @Accept      json
// @Produce     json
// @Success     200 {object} createInferenceJobResponse
// @Failure     500 {object} eResponse
// @Router      /inference-jobs [post]
// @Param       req body createInferenceJobRequest true "request"
func (r *InferenceJobController) create(c *gin.Context) {
	var req createInferenceJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - create")
		errorResponse(c, 400, "invalid request")

		return
	}

	job := entity.GenericJob{
		ID:     uuid.New().String(),
		Name:   "inference job",
		Status: "created",
	}

	entryPoint, err := r.u.CreateJob(job)
	if err != nil {
		r.l.Error(err, "http - v1 - create")
		errorResponse(c, 500, "database problems")
		return
	}

	go func() {
		time.Sleep(30 * time.Minute)
		err := r.u.DeleteJob(job.ID)
		if err != nil {
			r.l.Error(err, "http - v1 - create")
		}
	}()

	c.JSON(200, createInferenceJobResponse{JobId: job.ID, EntryPoint: entryPoint})
}

type getInferenceJobResponse struct {
	Job entity.GenericJob `json:"job"`
}

// @Summary     get inference job
// @Description get inference job
// @ID          get
// @Tags  	    inference-jobs
// @Accept      json
// @Produce     json
// @Param       id   path      string  true  "Job ID"
// @Success     200 {object} getInferenceJobResponse
// @Failure     500 {object} eResponse
// @Router      /inference-jobs/{id} [get]
func (r *InferenceJobController) get(c *gin.Context) {
	id := c.Param("id")

	job, err := r.u.GetJob(id)
	if err != nil {
		r.l.Error(err, "http - v1 - get")
		errorResponse(c, 500, "database problems")

		return
	}

	c.JSON(200, getInferenceJobResponse{job})
}

type listInferenceJobResponse struct {
	Jobs []entity.GenericJob `json:"jobs"`
}

// @Summary     list all inference jobs
// @Description list all inference jobs
// @ID          list
// @Tags  	    inference-jobs
// @Accept      json
// @Produce     json
// @Success     200 {object} listInferenceJobResponse
// @Failure     500 {object} eResponse
// @Router      /inference-jobs [get]
func (r *InferenceJobController) list(c *gin.Context) {
	jobs, err := r.u.GetAllJob()
	if err != nil {
		r.l.Error(err, "http - v1 - get")
		errorResponse(c, 500, "database problems")

		return
	}

	c.JSON(200, listInferenceJobResponse{jobs})
}

// @Summary     delete inference job
// @Description delete inference job
// @Tags  	    inference-jobs
// @Accept      json
// @Produce     json
// @Param       id   path      string  true  "Job ID"
// @Success     200 {object} sResponse
// @Failure     500 {object} eResponse
// @Router      /inference-jobs/{id} [delete]
func (r *InferenceJobController) delete(c *gin.Context) {
	id := c.Param("id")

	err := r.u.DeleteJob(id)
	if err != nil {
		r.l.Error(err, "http - v1 - delete")
		errorResponse(c, 500, "Internal problems, please try again later")

		return
	}

	successResponse(c, 200, "job deleted")
}
