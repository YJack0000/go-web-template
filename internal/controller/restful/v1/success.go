package v1

import (
	"github.com/gin-gonic/gin"
)

type sResponse struct {
	Success string `json:"success" example:"message"`
}

func successResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, sResponse{msg})
}
