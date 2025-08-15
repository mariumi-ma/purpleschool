package model

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

var dateFormat string = "2006-01-02 15:04:05"

type Product struct {
	gorm.Model
	Name        string
	Description string
	Images      pq.StringArray `gorm:"type:text[]"`
	Orders      []Order        `gorm:"many2many:order_products;"`
}

type Products []Product

type ProductResponse struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Images      []string `json:"images,omitempty"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

type ProductsResponse struct {
	Products []ProductResponse `json:"products"`
}

type UpdateProductRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
}

type CreateProductRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Images      []string `json:"images"`
}

func (p *Product) SetID(id int64) {
	p.ID = uint(id)
}

func (p *Product) ToResponse() *ProductResponse {
	return &ProductResponse{
		ID:          int64(p.ID),
		Name:        p.Name,
		Description: p.Description,
		Images:      p.Images,
		CreatedAt:   p.CreatedAt.Format(dateFormat),
		UpdatedAt:   p.UpdatedAt.Format(dateFormat),
	}
}

func (p *Products) ToResponse() ProductsResponse {
	result := ProductsResponse{
		Products: make([]ProductResponse, len(*p)),
	}

	for index, product := range *p {
		result.Products[index] = *product.ToResponse()
	}

	return result
}

func (p UpdateProductRequest) ToProduct() *Product {
	return &Product{
		Name:        p.Name,
		Description: p.Description,
		Images:      pq.StringArray(p.Images),
	}
}

func (p CreateProductRequest) ToProduct() *Product {
	return &Product{
		Name:        p.Name,
		Description: p.Description,
		Images:      pq.StringArray(p.Images),
	}
}
