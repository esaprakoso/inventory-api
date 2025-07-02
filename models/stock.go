package models

import (
	"time"

	"gorm.io/gorm"
)

type Stock struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Quantity  int            `json:"quantity"`
	ProductID uint           `json:"product_id"`
	Product   Product        `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"product"`
}
