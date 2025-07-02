package models

import "gorm.io/gorm"

type Stock struct {
	gorm.Model
	Quantity    int
	
	ProductID   uint
	Product     Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
