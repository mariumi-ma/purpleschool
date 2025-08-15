package order

import (
	"errors"

	"purpleschool/internal/database"
	"purpleschool/internal/model"

	"gorm.io/gorm"
)

var ErrOrderNotFound error = errors.New("order not found")

type OrderRepository struct {
	Database *database.Db
}

func NewOrderRepository(db *database.Db) *OrderRepository {
	return &OrderRepository{
		Database: db,
	}
}

func (r *OrderRepository) GetOrderByOrderIDAndUserID(orderID uint, userID uint) (*model.Order, error) {
	var order *model.Order

	result := r.Database.Preload("Products").
		Where("id = ?", orderID).
		Where("user_id = ?", userID).
		First(&order)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, result.Error
	}

	return order, nil
}

func (r *OrderRepository) GetOrdersByUserID(userID uint) (model.Orders, error) {
	var orders model.Orders

	tx := r.Database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := tx.Find(&orders).Where("user_id = ?", userID)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	// Ищем связи с продуктами
	if err := tx.Preload("Products").Find(&orders).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) CreateOrder(order *model.Order) (*model.Order, error) {

	tx := r.Database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Создаем сам заказ
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Сохраняем связи с продуктами (если они есть)
	if len(order.Products) > 0 {
		if err := tx.Model(order).Association("Products").Append(order.Products); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return order, nil
}
