package models

import "gorm.io/gorm"

// StockTransactionType mendefinisikan arah pergerakan stok (masuk/keluar)

type StockTransactionType string

const (
	StockTransactionTypeIn  StockTransactionType = "in"
	StockTransactionTypeOut StockTransactionType = "out"
)

// StockTransactionSubType mendefinisikan alasan spesifik dari transaksi
type StockTransactionSubType string

const (
	// Sub-tipe untuk stok masuk (In)
	SubTypePurchase   StockTransactionSubType = "purchase"   // Stok masuk dari pembelian ke supplier
	SubTypeReturn     StockTransactionSubType = "return"     // Stok masuk dari retur customer
	SubTypeTransferIn StockTransactionSubType = "transfer_in"// Stok masuk dari transfer gudang lain

	// Sub-tipe untuk stok keluar (Out)
	SubTypeSale        StockTransactionSubType = "sale"         // Stok keluar karena penjualan
	SubTypeDamaged     StockTransactionSubType = "damaged"      // Stok keluar karena rusak
	SubTypeExpired     StockTransactionSubType = "expired"      // Stok keluar karena kadaluarsa
	SubTypeTransferOut StockTransactionSubType = "transfer_out" // Stok keluar untuk transfer ke gudang lain

	// Sub-tipe umum
	SubTypeAdjustment StockTransactionSubType = "adjustment" // Penyesuaian hasil stock opname
)

type StockTransaction struct {
	gorm.Model
	StockID   uint                    `json:"stock_id"`
	Stock     Stock                   `json:"stock"`
	UserID    uint                    `json:"user_id"`
	User      User                    `json:"user"`
	Quantity  int                     `json:"quantity"` // Kuantitas selalu positif
	Type      StockTransactionType    `json:"type"`     // Tipe: 'in' atau 'out'
	SubType   StockTransactionSubType `json:"sub_type"` // Alasan spesifik transaksi
	Notes     string                  `json:"notes"`
}
