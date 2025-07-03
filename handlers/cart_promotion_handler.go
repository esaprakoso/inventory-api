package handlers

import (
	"net/http"

	"pos/database"
	"pos/models"

	"github.com/gin-gonic/gin"
)

// CreateCartPromotion handles the creation of a new cart promotion
func CreateCartPromotion(c *gin.Context) {
	var cartPromotion models.CartPromotion
	if err := c.ShouldBindJSON(&cartPromotion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body", "data": err.Error()})
		return
	}

	// Validate cart promotion type and data
	if cartPromotion.PromotionType != "percentage_discount" && cartPromotion.PromotionType != "fixed_discount" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid promotion type. Must be 'percentage_discount' or 'fixed_discount'"})
		return
	}

	if cartPromotion.DiscountValue <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "DiscountValue must be greater than 0"})
		return
	}

	if cartPromotion.MinimumPurchaseAmount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "MinimumPurchaseAmount must be greater than 0"})
		return
	}

	if cartPromotion.EndDate.Before(cartPromotion.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "End date cannot be before start date"})
		return
	}

	if err := database.DB.Create(&cartPromotion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not create cart promotion", "data": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Cart Promotion created", "data": cartPromotion})
}

// GetCartPromotions handles fetching all cart promotions
func GetCartPromotions(c *gin.Context) {
	var cartPromotions []models.CartPromotion
	database.DB.Find(&cartPromotions)
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Cart Promotions fetched", "data": cartPromotions})
}

// GetCartPromotion handles fetching a single cart promotion by ID
func GetCartPromotion(c *gin.Context) {
	id := c.Param("id")
	var cartPromotion models.CartPromotion
	if err := database.DB.First(&cartPromotion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Cart Promotion not found", "data": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Cart Promotion fetched", "data": cartPromotion})
}

// UpdateCartPromotion handles updating an existing cart promotion
func UpdateCartPromotion(c *gin.Context) {
	id := c.Param("id")
	var cartPromotion models.CartPromotion
	if err := c.ShouldBindJSON(&cartPromotion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body", "data": err.Error()})
		return
	}

	var existingCartPromotion models.CartPromotion
	if err := database.DB.First(&existingCartPromotion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Cart Promotion not found", "data": err.Error()})
		return
	}

	existingCartPromotion.PromotionType = cartPromotion.PromotionType
	existingCartPromotion.DiscountValue = cartPromotion.DiscountValue
	existingCartPromotion.MinimumPurchaseAmount = cartPromotion.MinimumPurchaseAmount
	existingCartPromotion.StartDate = cartPromotion.StartDate
	existingCartPromotion.EndDate = cartPromotion.EndDate

	// Validate cart promotion type and data
	if existingCartPromotion.PromotionType != "percentage_discount" && existingCartPromotion.PromotionType != "fixed_discount" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid promotion type. Must be 'percentage_discount' or 'fixed_discount'"})
		return
	}

	if existingCartPromotion.DiscountValue <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "DiscountValue must be greater than 0"})
		return
	}

	if existingCartPromotion.MinimumPurchaseAmount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "MinimumPurchaseAmount must be greater than 0"})
		return
	}

	if existingCartPromotion.EndDate.Before(existingCartPromotion.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "End date cannot be before start date"})
		return
	}

	if err := database.DB.Save(&existingCartPromotion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not update cart promotion", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Cart Promotion updated", "data": existingCartPromotion})
}

// DeleteCartPromotion handles deleting a cart promotion
func DeleteCartPromotion(c *gin.Context) {
	id := c.Param("id")
	var cartPromotion models.CartPromotion
	if err := database.DB.First(&cartPromotion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Cart Promotion not found", "data": err.Error()})
		return
	}

	if err := database.DB.Delete(&cartPromotion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not delete cart promotion", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Cart Promotion deleted"})
}
