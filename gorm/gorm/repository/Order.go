package repository

import (
	"fmt"

	"github.com/makoto-developer/golang_examples/gorm/gorm/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Get(orderID uint64) *model.Order
	ListByOrderID(orderIDs []uint64) ([]*model.Order, error)
	ListByUserID(userID uint64) ([]*model.Order, error)
	Create(order model.Order) error
	Update(order model.Order) error
	Delete(orderID uint64) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

// 注文IDで注文情報を取得
func (r *orderRepository) Get(orderID uint64) *model.Order {
	var order model.Order
	result := r.db.First(&order, orderID)
	if result.Error != nil {
		return nil
	}
	return &order
}

// 注文IDで注文を検索
func (r *orderRepository) ListByOrderID(orderIDs []uint64) ([]*model.Order, error) {
	var orders []*model.Order
	result := r.db.Where("id IN ?", orderIDs).Find(&orders)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list orders by order ids: %w", result.Error)
	}
	return orders, nil
}

// ユーザーIDで注文を全て取得
func (r *orderRepository) ListByUserID(userID uint64) ([]*model.Order, error) {
	var orders []*model.Order
	result := r.db.Where("user_id = ?", userID).Find(&orders)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list orders by user id: %w", result.Error)
	}
	return orders, nil
}

// 新規注文を作成
func (r *orderRepository) Create(order model.Order) error {
	result := r.db.Create(&order)
	if result.Error != nil {
		return fmt.Errorf("failed to create order: %w", result.Error)
	}
	return nil
}

// 注文を編集
func (r *orderRepository) Update(order model.Order) error {
	result := r.db.Save(&order)
	if result.Error != nil {
		return fmt.Errorf("failed to update order: %w", result.Error)
	}
	return nil
}

// 注文を削除(論理)
func (r *orderRepository) Delete(orderID uint64) error {
	result := r.db.Delete(&model.Order{}, orderID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete order: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("order not found with id: %d", orderID)
	}
	return nil
}
