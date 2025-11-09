// ============================================
// internal/repository/gorm_repository_test.go
// ============================================
package repository

import (
	"data-processing/internal/domain"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock
}

func TestNewGormRepository(t *testing.T) {
	db, _ := setupTestDB(t)

	repo := NewGormRepository(db)

	assert.NotNil(t, repo)
	assert.Implements(t, (*domain.ProductRepository)(nil), repo)
}

func TestGormRepository_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := setupTestDB(t)
		repo := NewGormRepository(db)

		product := &domain.Product{
			ID:        1,
			Name:      "Test Product",
			Brand:     "Test Brand",
			Category:  "Test Category",
			Price:     99.99,
			Currency:  "USD",
			Stock:     10,
			CreatedBy: "test_user",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "products"`)).
			WithArgs(
				sqlmock.AnyArg(), // Name
				sqlmock.AnyArg(), // Description
				sqlmock.AnyArg(), // Brand
				sqlmock.AnyArg(), // Category
				sqlmock.AnyArg(), // Price
				sqlmock.AnyArg(), // Currency
				sqlmock.AnyArg(), // Stock
				sqlmock.AnyArg(), // Ean
				sqlmock.AnyArg(), // Color
				sqlmock.AnyArg(), // Size
				sqlmock.AnyArg(), // Availability
				sqlmock.AnyArg(), // InternalId
				sqlmock.AnyArg(), // CreatedAt
				sqlmock.AnyArg(), // UpdatedAt
				sqlmock.AnyArg(), // CreatedBy
				sqlmock.AnyArg(), // ID
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := repo.Create(product)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error", func(t *testing.T) {
		db, mock := setupTestDB(t)
		repo := NewGormRepository(db)

		product := &domain.Product{
			ID:        1,
			Name:      "Test Product",
			Brand:     "Test Brand",
			Category:  "Test Category",
			Price:     99.99,
			CreatedBy: "test_user",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "products"`)).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err := repo.Create(product)

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGormRepository_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := setupTestDB(t)
		repo := NewGormRepository(db)

		product := &domain.Product{
			ID:        1,
			Name:      "Updated Product",
			Brand:     "Updated Brand",
			Category:  "Updated Category",
			Price:     149.99,
			Currency:  "USD",
			Stock:     20,
			CreatedBy: "test_user",
			UpdatedAt: time.Now(),
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(product)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error", func(t *testing.T) {
		db, mock := setupTestDB(t)
		repo := NewGormRepository(db)

		product := &domain.Product{
			ID:        1,
			Name:      "Updated Product",
			Brand:     "Updated Brand",
			Category:  "Updated Category",
			Price:     149.99,
			CreatedBy: "test_user",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products"`)).
			WillReturnError(errors.New("update failed"))
		mock.ExpectRollback()

		err := repo.Update(product)

		assert.Error(t, err)
		assert.Equal(t, "update failed", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGormRepository_FindById(t *testing.T) {
	t.Run("success - found", func(t *testing.T) {
		db, mock := setupTestDB(t)
		repo := NewGormRepository(db)

		expectedProduct := &domain.Product{
			ID:        1,
			Name:      "Test Product",
			Brand:     "Test Brand",
			Category:  "Test Category",
			Price:     99.99,
			Currency:  "USD",
			Stock:     10,
			CreatedBy: "test_user",
		}

		rows := sqlmock.NewRows([]string{
			"id", "name", "description", "brand", "category", "price",
			"currency", "stock", "ean", "color", "size", "availability",
			"internal_id", "created_at", "updated_at", "created_by",
		}).AddRow(
			expectedProduct.ID,
			expectedProduct.Name,
			expectedProduct.Description,
			expectedProduct.Brand,
			expectedProduct.Category,
			expectedProduct.Price,
			expectedProduct.Currency,
			expectedProduct.Stock,
			expectedProduct.Ean,
			expectedProduct.Color,
			expectedProduct.Size,
			expectedProduct.Availability,
			expectedProduct.InternalId,
			time.Now(),
			time.Now(),
			expectedProduct.CreatedBy,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE id = $1`)).
			WithArgs(1, 1).
			WillReturnRows(rows)

		product, err := repo.FindById(1)

		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, expectedProduct.ID, product.ID)
		assert.Equal(t, expectedProduct.Name, product.Name)
		assert.Equal(t, expectedProduct.Brand, product.Brand)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found - returns nil", func(t *testing.T) {
		db, mock := setupTestDB(t)
		repo := NewGormRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE id = $1`)).
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		product, err := repo.FindById(999)

		assert.NoError(t, err)
		assert.Nil(t, product)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock := setupTestDB(t)
		repo := NewGormRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE id = $1`)).
			WithArgs(1, 1).
			WillReturnError(errors.New("database connection error"))

		product, err := repo.FindById(1)

		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, "database connection error", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGormRepository_BulkUpsert(t *testing.T) {
	t.Run("success - with products", func(t *testing.T) {
		db, mock := setupTestDB(t)
		repo := NewGormRepository(db)

		products := []*domain.Product{
			{
				ID:        1,
				Name:      "Product 1",
				Brand:     "Brand 1",
				Category:  "Category 1",
				Price:     99.99,
				Currency:  "USD",
				Stock:     10,
				CreatedBy: "test_user",
			},
			{
				ID:        2,
				Name:      "Product 2",
				Brand:     "Brand 2",
				Category:  "Category 2",
				Price:     149.99,
				Currency:  "USD",
				Stock:     20,
				CreatedBy: "test_user",
			},
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "products"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectCommit()

		err := repo.BulkUpsert(products)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success - empty products", func(t *testing.T) {
		db, _ := setupTestDB(t)
		repo := NewGormRepository(db)

		products := []*domain.Product{}

		err := repo.BulkUpsert(products)

		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		db, mock := setupTestDB(t)
		repo := NewGormRepository(db)

		products := []*domain.Product{
			{
				ID:        1,
				Name:      "Product 1",
				Brand:     "Brand 1",
				Category:  "Category 1",
				Price:     99.99,
				CreatedBy: "test_user",
			},
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "products"`)).
			WillReturnError(errors.New("bulk insert failed"))
		mock.ExpectRollback()

		err := repo.BulkUpsert(products)

		assert.Error(t, err)
		assert.Equal(t, "bulk insert failed", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGormRepository_GetAll(t *testing.T) {
	t.Run("success - with products", func(t *testing.T) {
		db, mock := setupTestDB(t)
		repo := NewGormRepository(db)

		rows := sqlmock.NewRows([]string{
			"id", "name", "description", "brand", "category", "price",
			"currency", "stock", "ean", "color", "size", "availability",
			"internal_id", "created_at", "updated_at", "created_by",
		}).
			AddRow(1, "Product 1", "Desc 1", "Brand 1", "Category 1", 99.99,
				"USD", 10, "EAN1", "Red", "M", "In Stock", 101,
				time.Now(), time.Now(), "test_user").
			AddRow(2, "Product 2", "Desc 2", "Brand 2", "Category 2", 149.99,
				"USD", 20, "EAN2", "Blue", "L", "In Stock", 102,
				time.Now(), time.Now(), "test_user")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products"`)).
			WillReturnRows(rows)

		products, err := repo.GetAll()

		assert.NoError(t, err)
		assert.NotNil(t, products)
		assert.Len(t, products, 2)
		assert.Equal(t, "Product 1", products[0].Name)
		assert.Equal(t, "Product 2", products[1].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success - empty result", func(t *testing.T) {
		db, mock := setupTestDB(t)
		repo := NewGormRepository(db)

		rows := sqlmock.NewRows([]string{
			"id", "name", "description", "brand", "category", "price",
			"currency", "stock", "ean", "color", "size", "availability",
			"internal_id", "created_at", "updated_at", "created_by",
		})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products"`)).
			WillReturnRows(rows)

		products, err := repo.GetAll()

		assert.NoError(t, err)
		assert.NotNil(t, products)
		assert.Len(t, products, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error", func(t *testing.T) {
		db, mock := setupTestDB(t)
		repo := NewGormRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products"`)).
			WillReturnError(errors.New("query failed"))

		products, err := repo.GetAll()

		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Equal(t, "query failed", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
