package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID               int64          `gorm:"primaryKey" json:"id"`
	OrderItemGroupID int64          `json:"order_item_group_id"`
	UserID           int64          `json:"user_id" gorm:"index"`
	Amount           int64          `json:"amount"`
	AmountWithoutTax int64          `json:"amount_without_tax"`
	Tax              int64          `json:"tax"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at"`
}

func (Order) TableName() string {
	return "online_shop.\"order\""
}
