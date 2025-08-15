package model

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID   uint
	Products []Product `gorm:"many2many:order_products;"`
}

type Orders []Order

type OrderResponse struct {
	OrderID  uint              `json:"order_id"`
	Products []ProductResponse `json:"products"`
}

type OrdersResponse struct {
	Orders []OrderResponse `json:"orders"`
}

type CreateOrderRequest struct {
	ProductIDs []uint `json:"product_ids" validate:"required"`
}

type CreateOrderResponse struct {
	OrderID uint `json:"order_id"`
}

func (o *Order) SetUserID(userID uint) *Order {
	o.UserID = userID

	return o
}

func (o *Order) SetProductsIDs(productsIDs []uint) *Order {
	o.Products = make([]Product, len(productsIDs))

	for index, productID := range productsIDs {
		o.Products[index].ID = productID
	}

	return o
}

func (o *Order) ToResponse() OrderResponse {
	result := OrderResponse{
		OrderID: o.ID,
	}

	result.Products = make([]ProductResponse, len(o.Products))

	for index, product := range o.Products {
		result.Products[index] = *product.ToResponse()
	}

	return result
}

func (o *Orders) ToResponse() OrdersResponse {
	result := OrdersResponse{
		Orders: make([]OrderResponse, len(*o)),
	}

	for index, order := range *o {
		result.Orders[index] = order.ToResponse()
	}

	return result
}
