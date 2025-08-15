package handler

import (
	"net/http"

	"github.com/NRUG-SRE/slm-handson/backend/internal/interface/api/presenter"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	presenter.SuccessResponse(c, http.StatusOK, gin.H{
		"status":  "ok",
		"service": "slm-handson-api",
	})
}
