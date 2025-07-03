package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderItem struct {
	ID                uint           `gorm:"primarykey" json:"id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Quantity          int            `json:"quantity"`
	Price             float64        `json:"price"` // Original price of the product at the time of sale
	DiscountedPrice   float64        `json:"discounted_price"` // Price after item-specific discount
	ItemDiscount      float64        `gorm:"default:0" json:"item_discount"` // Discount amount for this item
	IsFreeItem        bool           `gorm:"default:false" json:"is_free_item"`
	OrderID           uint           `json:"order_id"`
	Order             Order          `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	ProductID         uint           `json:"product_id"`
	Product           Product        `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}
