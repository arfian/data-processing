// ============================================
// internal/domain/entity.go
// ============================================
package domain

import (
	"time"
)

// Product represents the domain model
type Product struct {
	ID          uint   `gorm:"primarykey"`
	SKU         string `gorm:"uniqueIndex;not null"`
	Name        string `gorm:"not null"`
	Description string
	Price       float64 `gorm:"not null"`
	Stock       int     `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CSVRecord represents raw CSV data
type CSVRecord struct {
	SKU         string
	Name        string
	Description string
	Price       string
	Stock       string
	RowNumber   int
}

// ProcessJob represents a job to be processed
type ProcessJob struct {
	Record   *CSVRecord
	FilePath string
}

// ProcessResult holds processing statistics
type ProcessResult struct {
	Product   *Product
	IsUpdate  bool
	Error     error
	RowNumber int
	FilePath  string
}

// ProgressUpdate represents real-time progress
type ProgressUpdate struct {
	FileName       string
	TotalRecords   int
	ProcessedCount int
	Percentage     float64
	Inserted       int
	Updated        int
	Failed         int
	Message        string
}

// FinalResult holds final processing statistics
type FinalResult struct {
	TotalRecords   int
	Inserted       int
	Updated        int
	Failed         int
	Errors         []string
	ProcessingTime time.Duration
	FileResults    map[string]*FileResult
}

// FileResult holds per-file statistics
type FileResult struct {
	TotalRecords int
	Inserted     int
	Updated      int
	Failed       int
	Errors       []string
}

// ProductRepository defines repository interface
type ProductRepository interface {
	Create(product *Product) error
	Update(product *Product) error
	FindBySKU(sku string) (*Product, error)
	BulkUpsert(products []*Product) error
	GetAll() ([]*Product, error)
}

// CSVProcessorUsecase defines usecase interface
type CSVProcessorUsecase interface {
	ProcessCSVFiles(filePaths []string, progressChan chan<- *ProgressUpdate) (*FinalResult, error)
}

// Logger defines logger interface
type Logger interface {
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Progress(filePath string, processed, total int, percentage float64)
}
