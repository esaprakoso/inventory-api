package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name       string         `json:"name"`
	SKU        string         `gorm:"index" json:"sku"`
	Price      float64        `json:"price"`
	CategoryID *uint          `json:"category_id"`
	Category   Category       `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Promotions []Promotion    `gorm:"foreignKey:ProductID" json:"promotions,omitempty"`
	Quantity         int            `json:"quantity"`
	ReservedQuantity int            `json:"reserved_quantity"`
}
