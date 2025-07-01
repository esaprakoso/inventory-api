package handlers

import (
	"errors"
	"fmt"
	"inventory/database"
	"inventory/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetStocks(c *gin.Context) {
	var stocks []models.Stock
	database.DB.Find(&stocks)
	c.JSON(http.StatusOK, stocks)
}

func GetStockByWarehouseID(c *gin.Context) {
	WarehouseID := c.Param("warehouse_id")

	var stocks models.Stock
	database.DB.Where("warehouse_id = ?", WarehouseID).Preload("Product").Find(&stocks)

	c.JSON(http.StatusOK, stocks)
}

func UpsertStock(c *gin.Context) {
	type UpsertStockInput struct {
		ProductID   uint                           `json:"product_id" binding:"required,exists=products-id"`
		WarehouseID uint                           `json:"warehouse_id" binding:"required,exists=warehouses-id"`
		Quantity    uint                           `json:"quantity" binding:"required,gt=0"`
		Type        models.StockTransactionType    `json:"type" binding:"required,oneof=in out"`
		SubType     models.StockTransactionSubType `json:"sub_type" binding:"required"`
		Notes       string                         `json:"notes"`
	}

	var data UpsertStockInput
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Get user ID from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found in context"})
		return
	}
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid user ID format"})
		return
	}

	// The transaction type is now explicitly provided in the request
	transactionType := data.Type

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		var stock models.Stock

		// Lock the stock record for update to prevent race conditions
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("warehouse_id = ? AND product_id = ?", data.WarehouseID, data.ProductID).
			First(&stock).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Stock does not exist, can only create it with an 'in' transaction
				if transactionType == models.StockTransactionTypeOut {
					return fmt.Errorf("stock not found, cannot perform 'out' transaction")
				}
				stock = models.Stock{
					ProductID:   data.ProductID,
					WarehouseID: data.WarehouseID,
					Quantity:    int(data.Quantity),
				}
				if err := tx.Create(&stock).Error; err != nil {
					return fmt.Errorf("failed to create stock: %w", err)
				}
			} else {
				// Another error occurred
				return fmt.Errorf("failed to find stock: %w", err)
			}
		} else {
			// Stock exists, update its quantity based on transaction type
			var newQuantity int
			if transactionType == models.StockTransactionTypeIn {
				newQuantity = stock.Quantity + int(data.Quantity)
			} else {
				newQuantity = stock.Quantity - int(data.Quantity)
			}

			if newQuantity < 0 {
				return fmt.Errorf("insufficient stock")
			}
			stock.Quantity = newQuantity
			if err := tx.Save(&stock).Error; err != nil {
				return fmt.Errorf("failed to update stock: %w", err)
			}
		}

		// Create the stock transaction log
		transaction := models.StockTransaction{
			StockID:  stock.ID,
			UserID:   uint(userID),
			Quantity: int(data.Quantity), // Log the positive quantity of the change
			Type:     transactionType,
			SubType:  data.SubType,
			Notes:    data.Notes,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return fmt.Errorf("failed to create transaction log: %w", err)
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock updated successfully"})
}
