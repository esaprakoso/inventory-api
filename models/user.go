package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Username     string         `gorm:"unique" json:"username"`
	Password     string         `json:"-"`
	Name         string         `json:"name"`
	RefreshToken string         `json:"-"`
	Role         string         `gorm:"default:'user'" json:"role"`
}