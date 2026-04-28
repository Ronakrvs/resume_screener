package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nvl258/resume-screener/internal/api/middleware"
	"github.com/nvl258/resume-screener/internal/model"
	"github.com/nvl258/resume-screener/internal/service"
)

type JobHandler struct {
	svc *service.JobService
}

func NewJobHandler(svc *service.JobService) *JobHandler {
	return &JobHandler{svc: svc}
}

func (h *JobHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req model.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	job, err := h.svc.Create(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, job)
}

func (h *JobHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	jobs, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"jobs": jobs})
}
