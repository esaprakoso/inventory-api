package models

import "gorm.io/gorm"

// StockTransactionType defines the direction of stock movement (in/out)

type StockTransactionType string

const (
	StockTransactionTypeIn  StockTransactionType = "in"
	StockTransactionTypeOut StockTransactionType = "out"
)

// StockTransactionSubType defines the specific reason for the transaction
type StockTransactionSubType string

const (
	// Sub-types for stock in (In)
	SubTypePurchase   StockTransactionSubType = "purchase"   // Stok masuk dari pembelian ke supplier
	SubTypeReturn     StockTransactionSubType = "return"     // Stok masuk dari retur customer
	SubTypeTransferIn StockTransactionSubType = "transfer_in"// Stok masuk dari transfer gudang lain

	// Sub-types for stock out (Out)
	SubTypeSale        StockTransactionSubType = "sale"         // Stok keluar karena penjualan
	SubTypeDamaged     StockTransactionSubType = "damaged"      // Stok keluar karena rusak
	SubTypeExpired     StockTransactionSubType = "expired"      // Stok keluar karena kadaluarsa
	SubTypeTransferOut StockTransactionSubType = "transfer_out" // Stok keluar untuk transfer ke gudang lain

	// General sub-type
	SubTypeAdjustment StockTransactionSubType = "adjustment" // Penyesuaian hasil stock opname
)

type StockTransaction struct {
	gorm.Model
	StockID   uint                    `json:"stock_id"`
	Stock     Stock                   `json:"stock"`
	UserID    uint                    `json:"user_id"`
	User      User                    `json:"user"`
	Quantity  int                     `json:"quantity"` // Quantity is always positive
	Type      StockTransactionType    `json:"type"`     // Type: 'in' or 'out'
	SubType   StockTransactionSubType `json:"sub_type"` // Specific reason for the transaction
	Notes     string                  `json:"notes"`
}
