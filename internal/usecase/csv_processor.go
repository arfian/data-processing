// ============================================
// internal/usecase/csv_processor.go
// ============================================
package usecase

import (
	"data-processing/internal/domain"
	"data-processing/pkg/csv"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type csvProcessorUsecase struct {
	repo        domain.ProductRepository
	logger      domain.Logger
	csvReader   *csv.Reader
	workerCount int
	batchSize   int
}

func NewCSVProcessorUsecase(
	repo domain.ProductRepository,
	logger domain.Logger,
	workerCount int,
	batchSize int,
) domain.CSVProcessorUsecase {
	return &csvProcessorUsecase{
		repo:        repo,
		logger:      logger,
		csvReader:   csv.NewReader(),
		workerCount: workerCount,
		batchSize:   batchSize,
	}
}

func (u *csvProcessorUsecase) ProcessCSVFiles(
	filePaths []string,
	progressChan chan<- *domain.ProgressUpdate,
) (*domain.FinalResult, error) {
	start := time.Now()
	u.logger.Info("Starting CSV processing with %d workers", u.workerCount)

	finalResult := &domain.FinalResult{
		FileResults: make(map[string]*domain.FileResult),
	}

	// Process each file
	for _, filePath := range filePaths {
		u.logger.Info("Processing file: %s", filePath)

		fileResult, err := u.processFileWithWorkers(filePath, progressChan)
		if err != nil {
			u.logger.Error("Failed to process file %s: %v", filePath, err)
			finalResult.Errors = append(finalResult.Errors, fmt.Sprintf("File %s: %v", filePath, err))
			continue
		}

		finalResult.FileResults[filePath] = fileResult
		finalResult.TotalRecords += fileResult.TotalRecords
		finalResult.Inserted += fileResult.Inserted
		finalResult.Updated += fileResult.Updated
		finalResult.Failed += fileResult.Failed
		finalResult.Errors = append(finalResult.Errors, fileResult.Errors...)
	}

	finalResult.ProcessingTime = time.Since(start)
	u.logger.Info("Processing completed in %v", finalResult.ProcessingTime)
	u.logger.Info("Total: %d | Inserted: %d | Updated: %d | Failed: %d",
		finalResult.TotalRecords, finalResult.Inserted, finalResult.Updated, finalResult.Failed)

	return finalResult, nil
}

func (u *csvProcessorUsecase) processFileWithWorkers(
	filePath string,
	progressChan chan<- *domain.ProgressUpdate,
) (*domain.FileResult, error) {
	// Read CSV file
	records, err := u.csvReader.ReadCSV(filePath)
	if err != nil {
		return nil, err
	}

	totalRecords := len(records)
	u.logger.Info("File %s: Found %d records", filePath, totalRecords)

	// Create channels
	jobChan := make(chan *domain.ProcessJob, totalRecords)
	resultChan := make(chan *domain.ProcessResult, totalRecords)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < u.workerCount; i++ {
		wg.Add(1)
		go u.worker(i+1, jobChan, resultChan, &wg)
	}

	// Send jobs to workers
	go func() {
		for _, record := range records {
			jobChan <- &domain.ProcessJob{
				Record:   record,
				FilePath: filePath,
			}
		}
		close(jobChan)
	}()

	// Close workers
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results and send progress updates
	fileResult := &domain.FileResult{
		TotalRecords: totalRecords,
	}

	processedCount := 0
	batch := make([]*domain.Product, 0, u.batchSize)

	for result := range resultChan {
		processedCount++

		if result.Error != nil {
			fileResult.Failed++
			errorMsg := fmt.Sprintf("Row %d (SKU: %s): %v",
				result.RowNumber, result.Product.SKU, result.Error)
			fileResult.Errors = append(fileResult.Errors, errorMsg)
			u.logger.Error(errorMsg)
		} else {
			if result.IsUpdate {
				fileResult.Updated++
			} else {
				fileResult.Inserted++
			}
			batch = append(batch, result.Product)

			// Batch upsert
			if len(batch) >= u.batchSize {
				if err := u.repo.BulkUpsert(batch); err != nil {
					u.logger.Error("Batch upsert failed: %v", err)
				}
				batch = batch[:0]
			}
		}

		// Send progress update every 10 records or at completion
		if processedCount%10 == 0 || processedCount == totalRecords {
			percentage := float64(processedCount) / float64(totalRecords) * 100
			u.logger.Progress(filePath, processedCount, totalRecords, percentage)

			if progressChan != nil {
				progressChan <- &domain.ProgressUpdate{
					FileName:       filePath,
					TotalRecords:   totalRecords,
					ProcessedCount: processedCount,
					Percentage:     percentage,
					Inserted:       fileResult.Inserted,
					Updated:        fileResult.Updated,
					Failed:         fileResult.Failed,
					Message:        fmt.Sprintf("Processing %s: %.2f%% complete", filePath, percentage),
				}
			}
		}
	}

	// Final batch upsert
	if len(batch) > 0 {
		if err := u.repo.BulkUpsert(batch); err != nil {
			u.logger.Error("Final batch upsert failed: %v", err)
		}
	}

	return fileResult, nil
}

func (u *csvProcessorUsecase) worker(
	id int,
	jobs <-chan *domain.ProcessJob,
	results chan<- *domain.ProcessResult,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for job := range jobs {
		result := u.processRecord(job)
		results <- result
	}
}

func (u *csvProcessorUsecase) processRecord(job *domain.ProcessJob) *domain.ProcessResult {
	record := job.Record

	// Convert CSV record to Product
	product, err := u.convertToProduct(record)
	if err != nil {
		return &domain.ProcessResult{
			Product:   product,
			Error:     err,
			RowNumber: record.RowNumber,
			FilePath:  job.FilePath,
		}
	}

	// Check if product exists
	existing, err := u.repo.FindBySKU(record.SKU)
	if err != nil {
		return &domain.ProcessResult{
			Product:   product,
			Error:     err,
			RowNumber: record.RowNumber,
			FilePath:  job.FilePath,
		}
	}

	isUpdate := false
	if existing != nil {
		product.ID = existing.ID
		product.CreatedAt = existing.CreatedAt
		isUpdate = true
	}

	return &domain.ProcessResult{
		Product:   product,
		IsUpdate:  isUpdate,
		RowNumber: record.RowNumber,
		FilePath:  job.FilePath,
	}
}

func (u *csvProcessorUsecase) convertToProduct(record *domain.CSVRecord) (*domain.Product, error) {
	price, err := strconv.ParseFloat(record.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid price: %v", err)
	}

	stock, err := strconv.Atoi(record.Stock)
	if err != nil {
		return nil, fmt.Errorf("invalid stock: %v", err)
	}

	if record.SKU == "" || record.Name == "" {
		return nil, fmt.Errorf("SKU and Name are required")
	}

	return &domain.Product{
		SKU:         record.SKU,
		Name:        record.Name,
		Description: record.Description,
		Price:       price,
		Stock:       stock,
	}, nil
}
