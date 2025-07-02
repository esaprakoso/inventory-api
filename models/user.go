package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"unique"`
	Password     string `json:"-"`
	Name         string
	RefreshToken string      `json:"-"`
	Role         string      `gorm:"default:'user'"`
	
}
