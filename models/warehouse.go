package models

import "gorm.io/gorm"

type Warehouse struct {
	gorm.Model
	Name     string
	Location string
	OwnerID  uint
	Owner    User `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Stocks   []Stock `gorm:"foreignKey:WarehouseID" json:"stocks"`
}
