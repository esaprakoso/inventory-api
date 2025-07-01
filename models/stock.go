package models

import "gorm.io/gorm"

type Stock struct {
	gorm.Model
	Quantity    int
	WarehouseID uint
	Warehouse   Warehouse `gorm:"foreignKey:WarehouseID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ProductID   uint
	Product     Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
