package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name       string
	SKU        string `gorm:"index"`
	Price      float64
	CategoryID *uint
	Category   Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Stocks     []Stock  `gorm:"foreignKey:ProductID" json:"-"`
}
