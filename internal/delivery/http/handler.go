// ============================================
// internal/delivery/http/handler.go
// ============================================
package handler

import (
	"net/http"

	"data-processing/internal/domain"

	"github.com/gin-gonic/gin"

	_ "data-processing/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	usecase domain.CSVProcessorUsecase
}

func NewHandler(usecase domain.CSVProcessorUsecase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := r.Group("/api/v1")
	{
		api.POST("/csv/process", h.ProcessCSV)
	}
}

type ProcessCSVRequest struct {
	FilePaths []string `json:"file_paths" binding:"required"`
}

// @BasePath /api/v1

// @Summary Process CSV
// @Description Insert / Update Process CSV
// @Tags csv
// @Accept json
// @Produce json
// @Param csv body ProcessCSVRequest true "Array Path CSV"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /csv/process [post]
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
