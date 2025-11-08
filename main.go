// ============================================
// cmd/api/main.go
// ============================================
package main

import (
	"log"

	"data-processing/config"
	handler "data-processing/internal/delivery/http"
	"data-processing/internal/repository"
	"data-processing/internal/usecase"
	"data-processing/pkg/database"
	"data-processing/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	appLogger := logger.NewLogger()

	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	repo := repository.NewGormRepository(db)
	uc := usecase.NewCSVProcessorUsecase(repo, appLogger, cfg.WorkerCount, cfg.BatchSize)
	handler := handler.NewHandler(uc)

	r := gin.Default()
	handler.RegisterRoutes(r)

	log.Printf("Server starting on port %s with %d workers", cfg.ServerPort, cfg.WorkerCount)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
