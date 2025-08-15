package product

import (
	"errors"

	"purpleschool/internal/database"
	"purpleschool/internal/model"

	"gorm.io/gorm"
)

var ErrProductNotFound error = errors.New("product not found")

type ProductRepository struct {
	Database *database.Db
}

func NewProductRepository(db *database.Db) *ProductRepository {
	return &ProductRepository{
		Database: db,
	}
}

func (r *ProductRepository) GetProductByID(id uint) (*model.Product, error) {
	var product *model.Product

	result := r.Database.First(&product, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, result.Error
	}

	return product, nil
}

func (r *ProductRepository) GetProducts() (model.Products, error) {
	var products model.Products

	result := r.Database.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func (r *ProductRepository) CreateProduct(product *model.Product) (*model.Product, error) {
	result := r.Database.Create(product)
	if result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (r *ProductRepository) UpdateProduct(product *model.Product) (*model.Product, error) {
	result := r.Database.Updates(product)
	if result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (r *ProductRepository) DeleteProduct(id uint) error {
	result := r.Database.Delete(&model.Product{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
