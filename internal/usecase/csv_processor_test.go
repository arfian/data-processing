// ============================================
// internal/usecase/csv_processor_test.go
// ============================================
package usecase

import (
	"data-processing/internal/domain"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCSVProcessorUsecase(t *testing.T) {
	mockRepo := domain.NewMockProductRepository(t)
	mockLogger := domain.NewMockLogger(t)

	usecase := NewCSVProcessorUsecase(mockRepo, mockLogger, 4, 100)

	assert.NotNil(t, usecase)
	assert.Implements(t, (*domain.CSVProcessorUsecase)(nil), usecase)
}

func TestConvertToProduct(t *testing.T) {
	mockRepo := domain.NewMockProductRepository(t)
	mockLogger := domain.NewMockLogger(t)
	u := &csvProcessorUsecase{
		repo:   mockRepo,
		logger: mockLogger,
	}
	_ = u // use variable to avoid unused warning

	t.Run("success - valid record", func(t *testing.T) {
		record := &domain.CSVRecord{
			ID:           "1",
			Name:         "Test Product",
			Description:  "Test Description",
			Brand:        "Test Brand",
			Category:     "Test Category",
			Price:        "99.99",
			Currency:     "USD",
			Stock:        "10",
			Ean:          "1234567890",
			Color:        "Red",
			Size:         "M",
			Availability: "In Stock",
			InternalId:   "100",
			RowNumber:    1,
		}

		product, err := u.convertToProduct(record)

		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, 1, product.ID)
		assert.Equal(t, "Test Product", product.Name)
		assert.Equal(t, "Test Description", product.Description)
		assert.Equal(t, "Test Brand", product.Brand)
		assert.Equal(t, "Test Category", product.Category)
		assert.Equal(t, 99.99, product.Price)
		assert.Equal(t, "USD", product.Currency)
		assert.Equal(t, 10, product.Stock)
		assert.Equal(t, "1234567890", product.Ean)
		assert.Equal(t, "Red", product.Color)
		assert.Equal(t, "M", product.Size)
		assert.Equal(t, "In Stock", product.Availability)
		assert.Equal(t, 100, product.InternalId)
		assert.Equal(t, "system", product.CreatedBy)
	})

	t.Run("error - invalid price", func(t *testing.T) {
		record := &domain.CSVRecord{
			ID:         "1",
			Name:       "Test Product",
			Brand:      "Test Brand",
			Category:   "Test Category",
			Price:      "invalid",
			Stock:      "10",
			InternalId: "100",
		}

		product, err := u.convertToProduct(record)

		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Contains(t, err.Error(), "invalid price")
	})

	t.Run("error - invalid stock", func(t *testing.T) {
		record := &domain.CSVRecord{
			ID:         "1",
			Name:       "Test Product",
			Brand:      "Test Brand",
			Category:   "Test Category",
			Price:      "99.99",
			Stock:      "invalid",
			InternalId: "100",
		}

		product, err := u.convertToProduct(record)

		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Contains(t, err.Error(), "invalid stock")
	})

	t.Run("error - invalid id", func(t *testing.T) {
		record := &domain.CSVRecord{
			ID:         "invalid",
			Name:       "Test Product",
			Brand:      "Test Brand",
			Category:   "Test Category",
			Price:      "99.99",
			Stock:      "10",
			InternalId: "100",
		}

		product, err := u.convertToProduct(record)

		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Contains(t, err.Error(), "invalid stock")
	})

	t.Run("error - invalid internal id", func(t *testing.T) {
		record := &domain.CSVRecord{
			ID:         "1",
			Name:       "Test Product",
			Brand:      "Test Brand",
			Category:   "Test Category",
			Price:      "99.99",
			Stock:      "10",
			InternalId: "invalid",
		}

		product, err := u.convertToProduct(record)

		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Contains(t, err.Error(), "invalid stock")
	})

	t.Run("error - empty name", func(t *testing.T) {
		record := &domain.CSVRecord{
			ID:         "1",
			Name:       "",
			Brand:      "Test Brand",
			Category:   "Test Category",
			Price:      "99.99",
			Stock:      "10",
			InternalId: "100",
		}

		product, err := u.convertToProduct(record)

		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Contains(t, err.Error(), "required")
	})
}

func TestProcessRecord(t *testing.T) {
	t.Run("success - new product", func(t *testing.T) {
		mockRepo := domain.NewMockProductRepository(t)
		mockLogger := domain.NewMockLogger(t)
		u := &csvProcessorUsecase{
			repo:   mockRepo,
			logger: mockLogger,
		}

		job := &domain.ProcessJob{
			Record: &domain.CSVRecord{
				ID:         "1",
				Name:       "Test Product",
				Brand:      "Test Brand",
				Category:   "Test Category",
				Price:      "99.99",
				Stock:      "10",
				InternalId: "100",
				RowNumber:  1,
			},
			FilePath: "/test/file.csv",
		}

		mockRepo.On("FindById", 1).Return(nil, nil)

		result := u.processRecord(job)

		assert.NotNil(t, result)
		assert.NoError(t, result.Error)
		assert.False(t, result.IsUpdate)
		assert.Equal(t, 1, result.RowNumber)
		assert.Equal(t, "/test/file.csv", result.FilePath)
		assert.NotNil(t, result.Product)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success - existing product (update)", func(t *testing.T) {
		mockRepo := domain.NewMockProductRepository(t)
		mockLogger := domain.NewMockLogger(t)
		u := &csvProcessorUsecase{
			repo:   mockRepo,
			logger: mockLogger,
		}

		job := &domain.ProcessJob{
			Record: &domain.CSVRecord{
				ID:         "1",
				Name:       "Updated Product",
				Brand:      "Test Brand",
				Category:   "Test Category",
				Price:      "149.99",
				Stock:      "20",
				InternalId: "100",
				RowNumber:  2,
			},
			FilePath: "/test/file.csv",
		}

		existingProduct := &domain.Product{
			ID:        1,
			Name:      "Old Product",
			Brand:     "Test Brand",
			Category:  "Test Category",
			Price:     99.99,
			Stock:     10,
			CreatedBy: "system",
		}

		mockRepo.On("FindById", 1).Return(existingProduct, nil)

		result := u.processRecord(job)

		assert.NotNil(t, result)
		assert.NoError(t, result.Error)
		assert.True(t, result.IsUpdate)
		assert.Equal(t, 2, result.RowNumber)
		assert.Equal(t, "/test/file.csv", result.FilePath)
		assert.NotNil(t, result.Product)
		assert.Equal(t, existingProduct.ID, result.Product.ID)
		assert.Equal(t, existingProduct.CreatedAt, result.Product.CreatedAt)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - invalid conversion", func(t *testing.T) {
		mockRepo := domain.NewMockProductRepository(t)
		mockLogger := domain.NewMockLogger(t)
		u := &csvProcessorUsecase{
			repo:   mockRepo,
			logger: mockLogger,
		}

		job := &domain.ProcessJob{
			Record: &domain.CSVRecord{
				ID:         "1",
				Name:       "Test Product",
				Brand:      "Test Brand",
				Category:   "Test Category",
				Price:      "invalid",
				Stock:      "10",
				InternalId: "100",
				RowNumber:  3,
			},
			FilePath: "/test/file.csv",
		}

		result := u.processRecord(job)

		assert.NotNil(t, result)
		assert.Error(t, result.Error)
		assert.Contains(t, result.Error.Error(), "invalid price")
		assert.Equal(t, 3, result.RowNumber)
	})

	t.Run("error - repository error", func(t *testing.T) {
		mockRepo := domain.NewMockProductRepository(t)
		mockLogger := domain.NewMockLogger(t)
		u := &csvProcessorUsecase{
			repo:   mockRepo,
			logger: mockLogger,
		}

		job := &domain.ProcessJob{
			Record: &domain.CSVRecord{
				ID:         "1",
				Name:       "Test Product",
				Brand:      "Test Brand",
				Category:   "Test Category",
				Price:      "99.99",
				Stock:      "10",
				InternalId: "100",
				RowNumber:  4,
			},
			FilePath: "/test/file.csv",
		}

		mockRepo.On("FindById", 1).Return(nil, errors.New("database error"))

		result := u.processRecord(job)

		assert.NotNil(t, result)
		assert.Error(t, result.Error)
		assert.Equal(t, "database error", result.Error.Error())
		assert.Equal(t, 4, result.RowNumber)
		mockRepo.AssertExpectations(t)
	})
}

func TestWorker(t *testing.T) {
	mockRepo := domain.NewMockProductRepository(t)
	mockLogger := domain.NewMockLogger(t)
	u := &csvProcessorUsecase{
		repo:   mockRepo,
		logger: mockLogger,
	}

	jobs := make(chan *domain.ProcessJob, 2)
	results := make(chan *domain.ProcessResult, 2)

	job1 := &domain.ProcessJob{
		Record: &domain.CSVRecord{
			ID:         "1",
			Name:       "Product 1",
			Brand:      "Brand 1",
			Category:   "Category 1",
			Price:      "99.99",
			Stock:      "10",
			InternalId: "100",
			RowNumber:  1,
		},
		FilePath: "/test/file.csv",
	}

	job2 := &domain.ProcessJob{
		Record: &domain.CSVRecord{
			ID:         "2",
			Name:       "Product 2",
			Brand:      "Brand 2",
			Category:   "Category 2",
			Price:      "149.99",
			Stock:      "20",
			InternalId: "200",
			RowNumber:  2,
		},
		FilePath: "/test/file.csv",
	}

	mockRepo.On("FindById", 1).Return(nil, nil)
	mockRepo.On("FindById", 2).Return(nil, nil)

	jobs <- job1
	jobs <- job2
	close(jobs)

	var wg sync.WaitGroup
	wg.Add(1)
	go u.worker(1, jobs, results, &wg)
	wg.Wait()
	close(results)

	result1 := <-results
	result2 := <-results

	assert.NotNil(t, result1)
	assert.NotNil(t, result2)
	assert.NoError(t, result1.Error)
	assert.NoError(t, result2.Error)
	mockRepo.AssertExpectations(t)
}
