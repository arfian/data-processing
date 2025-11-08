// ============================================
// pkg/logger/logger.go
// ============================================
package logger

import (
	"log"
	"time"
)

type AppLogger struct{}

func NewLogger() *AppLogger {
	return &AppLogger{}
}

func (l *AppLogger) Info(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

func (l *AppLogger) Error(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

func (l *AppLogger) Debug(format string, args ...interface{}) {
	log.Printf("[DEBUG] "+format, args...)
}

func (l *AppLogger) Progress(filePath string, processed, total int, percentage float64) {
	timestamp := time.Now().Format("15:04:05")
	log.Printf("[PROGRESS] [%s] %s: %d/%d (%.2f%%)",
		timestamp, filePath, processed, total, percentage)
}
