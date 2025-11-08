// ============================================
// internal/delivery/http/handler.go
// ============================================
package handler

import (
	"net/http"

	"data-processing/internal/domain"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase domain.CSVProcessorUsecase
}

func NewHandler(usecase domain.CSVProcessorUsecase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.POST("/csv/process", h.ProcessCSV)
	}
}

type ProcessCSVRequest struct {
	FilePaths []string `json:"file_paths" binding:"required"`
}

func (h *Handler) ProcessCSV(c *gin.Context) {
	var req ProcessCSVRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Channel for progress updates (optional for WebSocket in future)
	progressChan := make(chan *domain.ProgressUpdate, 100)

	// Close progress channel when done
	defer close(progressChan)

	// Process in background and collect progress
	go func() {
		for progress := range progressChan {
			// In production, send via WebSocket or SSE
			_ = progress
		}
	}()

	result, err := h.usecase.ProcessCSVFiles(req.FilePaths, progressChan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "CSV files processed successfully",
		"result":  result,
	})
}
