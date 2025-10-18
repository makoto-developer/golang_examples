package model

type Order struct {
	ID               int64 `gorm:"primaryKey" json:"id"`
	OrderItemGroupID int64 `json:"order_item_group_id"`
	UserID           int64 `json:"user_id" gorm:"index"`
	Amount           int64 `json:"amount"`
	AmountWithoutTax int64 `json:"amount_without_tax"`
	Tax              int64 `json:"tax"`
}

func (Order) TableName() string {
	return "online_shop.\"order\""
}
