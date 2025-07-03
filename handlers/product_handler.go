package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"pos/database"
	"pos/models"

	"github.com/gin-gonic/gin"
	gorm "gorm.io/gorm"
	"gorm.io/gorm/clause"

	"pos/utils"
)

type ProductResponse struct {
	models.Product
	CategoryName     string                   `json:"category_name"`
	Quantity         int                      `json:"quantity"`
	ReservedQuantity int                      `json:"reserved_quantity"`
	DiscountedPrice  float64                  `json:"discounted_price"`
	ActivePromotion  *models.ProductPromotion `json:"active_promotion,omitempty"`
}

func GetAllProducts(c *gin.Context) {
	var products []models.Product
	var total int64

	page, _ := utils.GetInt(c.DefaultQuery("page", "1"))
	limit, _ := utils.GetInt(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	query := database.DB.Model(&models.Product{}).Preload("Category").Preload("Promotions")

	query.Count(&total)
	query.Limit(limit).Offset(offset).Find(&products)

	var productResponses []ProductResponse
	for _, p := range products {
		discountedPrice, activePromotion := utils.CalculateDiscountedPrice(p)
		productResponses = append(productResponses, ProductResponse{
			Product:         p,
			CategoryName:    p.Category.Name,
			Quantity:        p.Quantity,
			DiscountedPrice: discountedPrice,
			ActivePromotion: activePromotion,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  productResponses,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func StoreProduct(c *gin.Context) {
	type CreateProductInput struct {
		Name       string  `json:"name" binding:"required"`
		Price      float64 `json:"price" binding:"required"`
		SKU        string  `json:"sku" binding:"required"`
		CategoryID *uint   `json:"category_id"`
	}

	var data CreateProductInput

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	isDup, err := utils.IsDuplicate[models.Product](database.DB, "SKU", data.SKU, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	if isDup {
		c.JSON(406, gin.H{"message": "Product SKU already exists"})
		return
	}

	product := models.Product{
		Name:             data.Name,
		SKU:              data.SKU,
		Price:            data.Price,
		CategoryID:       data.CategoryID,
		Quantity:         0, // Initial quantity
		ReservedQuantity: 0,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create product", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product Created"})
}

func GetProductByID(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	database.DB.Preload("Category").Preload("Promotions").First(&product, id)

	if product.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}

	discountedPrice, activePromotion := utils.CalculateDiscountedPrice(product)

	productResponse := ProductResponse{
		Product:          product,
		CategoryName:     product.Category.Name,
		Quantity:         product.Quantity,
		ReservedQuantity: product.ReservedQuantity,
		DiscountedPrice:  discountedPrice,
		ActivePromotion:  activePromotion,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": productResponse,
	})
}

func UpdateProductByID(c *gin.Context) {
	id := c.Param("id")

	type UpdateProductInput struct {
		Name       string  `json:"name" binding:"required"`
		Price      float64 `json:"price" binding:"required"`
		SKU        string  `json:"sku" binding:"required"`
		CategoryID *uint   `json:"category_id"`
	}
	var data UpdateProductInput

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var product models.Product
	database.DB.Preload("Category").First(&product, id)

	isDup, err := utils.IsDuplicate[models.Product](database.DB, "SKU", data.SKU, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	if isDup {
		c.JSON(406, gin.H{"message": "Product SKU already exists"})
		return
	}

	if product.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}

	product.Name = data.Name
	product.SKU = data.SKU
	product.Price = data.Price
	product.CategoryID = data.CategoryID

	database.DB.Save(&product)

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated",
	})
}

func DeleteProductByID(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	database.DB.First(&product, id)

	if product.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}

	database.DB.Delete(&product, id)
	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted",
	})
}

func UpdateProductStock(c *gin.Context) {
	id := c.Param("id")

	type UpdateStockInput struct {
		Quantity int                            `json:"quantity" binding:"required,gt=0"`
		Type     models.StockTransactionType    `json:"type" binding:"required,oneof=in out"`
		SubType  models.StockTransactionSubType `json:"sub_type" binding:"required"`
		Notes    string                         `json:"notes"`
	}

	var data UpdateStockInput
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Get user ID from context
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid User ID format in context"})
		return
	}

	transactionErr := database.DB.Transaction(func(tx *gorm.DB) error {
		var product models.Product

		// Lock the product record for update to prevent race conditions
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&product, id).Error

		if err != nil {
			return fmt.Errorf("product not found for ID %s: %w", id, err)
		}

		// Update product quantity based on transaction type
		var newQuantity int
		if data.Type == models.StockTransactionTypeIn {
			newQuantity = product.Quantity + data.Quantity
		} else {
			newQuantity = product.Quantity - data.Quantity
		}

		if newQuantity < 0 {
			return fmt.Errorf("insufficient stock: cannot go below zero")
		}
		product.Quantity = newQuantity
		if err := tx.Save(&product).Error; err != nil {
			return fmt.Errorf("failed to update product quantity: %w", err)
		}

		// Create the stock transaction log
		transaction := models.StockTransaction{
			ProductID: product.ID,
			UserID:    uint(userID),
			Quantity:  data.Quantity, // Log the positive quantity of the change
			Type:      data.Type,
			SubType:   data.SubType,
			Notes:     data.Notes,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return fmt.Errorf("failed to create transaction log: %w", err)
		}

		return nil
	})

	if transactionErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": transactionErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock updated successfully"})
}
