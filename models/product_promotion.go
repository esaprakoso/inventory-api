package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductPromotion struct {
	ID               uint           `gorm:"primarykey" json:"id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string"`
	ProductID        uint           `json:"product_id"`
	Product          Product        `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	PromotionType    string         `json:"promotion_type"` // e.g., "percentage_discount", "fixed_discount", "buy_x_get_y", "bundle_price"
	DiscountValue    float64        `json:"discount_value,omitempty"`
	BuyProductID     *uint          `json:"buy_product_id,omitempty"`
	GetProductID     *uint          `json:"get_product_id,omitempty"`
	RequiredQuantity *int           `json:"required_quantity,omitempty"` // For "bundle_price" or "buy_x_get_y"
	PromoPrice       *float64       `json:"promo_price,omitempty"`       // For "bundle_price"
	StartDate        time.Time      `json:"start_date"`
	EndDate          time.Time      `json:"end_date"`
}
