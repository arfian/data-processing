// ============================================
// internal/repository/gorm_repository.go
// ============================================
package repository

import (
	"data-processing/internal/domain"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) domain.ProductRepository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *gormRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *gormRepository) FindBySKU(sku string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Where("sku = ?", sku).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *gormRepository) BulkUpsert(products []*domain.Product) error {
	if len(products) == 0 {
		return nil
	}

	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "sku"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "description", "price", "stock", "updated_at"}),
	}).CreateInBatches(&products, 100).Error
}

func (r *gormRepository) GetAll() ([]*domain.Product, error) {
	var products []*domain.Product
	err := r.db.Find(&products).Error
	return products, err
}
