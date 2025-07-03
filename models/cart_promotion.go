package models

import (
	"time"

	"gorm.io/gorm"
)

type CartPromotion struct {
	ID                    uint           `gorm:"primarykey" json:"id"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	PromotionType         string         `json:"promotion_type"` // e.g., "percentage_discount", "fixed_discount"
	DiscountValue         float64        `json:"discount_value"`
	MinimumPurchaseAmount float64        `json:"minimum_purchase_amount"`
	StartDate             time.Time      `json:"start_date"`
	EndDate               time.Time      `json:"end_date"`
}
