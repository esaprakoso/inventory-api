package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name     string    `gorm:"index"`
	Products []Product `gorm:"foreignKey:CategoryID" json:"products"`
}
