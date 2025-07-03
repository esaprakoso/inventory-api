package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID                uint           `gorm:"primarykey" json:"id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	GrossTotal        float64        `json:"gross_total"`                          // GrossTotal adalah total harga dari semua item sebelum diskon per item diterapkan
	SubTotal          float64        `json:"sub_total"`                            // SubTotal adalah total harga dari semua item setelah diskon per item diterapkan
	ItemDiscountTotal float64        `gorm:"default:0" json:"item_discount_total"` // ItemDiscountTotal adalah total akumulasi diskon yang diberikan per item
	CartDiscount      float64        `gorm:"default:0" json:"cart_discount"`       // CartDiscount adalah diskon yang diterapkan pada total belanja (misal: diskon minimal, kupon)
	TotalAmount       float64        `json:"total_amount"`                         // TotalAmount adalah jumlah akhir yang harus dibayar pelanggan (SubTotal - CartDiscount)
	PaymentMethod     string         `json:"payment_method"`
	UserID            uint           `json:"user_id"`
	User              User           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	OrderItems        []OrderItem    `gorm:"foreignKey:OrderID" json:"order_items"`
}
