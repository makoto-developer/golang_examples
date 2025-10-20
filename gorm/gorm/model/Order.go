package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID               int64          `gorm:"primaryKey" json:"id"`
	OrderItemGroupID int64          `gorm:"order_item_group_id"`
	UserID           int64          `gorm:"user_id" gorm:"index"`
	Amount           int64          `gorm:"amount"`
	AmountWithoutTax int64          `gorm:"amount_without_tax"`
	Tax              int64          `gorm:"tax"`
	CreatedAt        time.Time      `gorm:"created_at"`
	UpdatedAt        time.Time      `gorm:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"deleted_at"`
}
