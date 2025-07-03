package handlers

import (
	"net/http"

	"pos/database"
	"pos/models"

	"github.com/gin-gonic/gin"
)

// CreatePromotion handles the creation of a new promotion
func CreatePromotion(c *gin.Context) {
	var promotion models.Promotion
	if err := c.ShouldBindJSON(&promotion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body", "data": err.Error()})
		return
	}

	// Validate promotion type and data
	if promotion.PromotionType == "buy_x_get_y" {
		if promotion.BuyProductID == nil || promotion.GetProductID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "BuyProductID and GetProductID are required for buy_x_get_y promotion"})
			return
		}
	} else if promotion.PromotionType == "percentage_discount" || promotion.PromotionType == "fixed_discount" {
		if promotion.DiscountValue <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "DiscountValue must be greater than 0 for discount promotions"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid promotion type"})
		return
	}

	if promotion.EndDate.Before(promotion.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "End date cannot be before start date"})
		return
	}

	if err := database.DB.Create(&promotion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not create promotion", "data": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Promotion created", "data": promotion})
}

// GetPromotions handles fetching all promotions
func GetPromotions(c *gin.Context) {
	var promotions []models.Promotion
	database.DB.Find(&promotions)
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Promotions fetched", "data": promotions})
}

// GetPromotion handles fetching a single promotion by ID
func GetPromotion(c *gin.Context) {
	id := c.Param("id")
	var promotion models.Promotion
	if err := database.DB.First(&promotion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Promotion not found", "data": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Promotion fetched", "data": promotion})
}

// UpdatePromotion handles updating an existing promotion
func UpdatePromotion(c *gin.Context) {
	id := c.Param("id")
	var promotion models.Promotion
	if err := c.ShouldBindJSON(&promotion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body", "data": err.Error()})
		return
	}

	var existingPromotion models.Promotion
	if err := database.DB.First(&existingPromotion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Promotion not found", "data": err.Error()})
		return
	}

	// Update fields
	existingPromotion.ProductID = promotion.ProductID
	existingPromotion.PromotionType = promotion.PromotionType
	existingPromotion.DiscountValue = promotion.DiscountValue
	existingPromotion.BuyProductID = promotion.BuyProductID
	existingPromotion.GetProductID = promotion.GetProductID
	existingPromotion.StartDate = promotion.StartDate
	existingPromotion.EndDate = promotion.EndDate

	// Validate promotion type and data
	if existingPromotion.PromotionType == "buy_x_get_y" {
		if existingPromotion.BuyProductID == nil || existingPromotion.GetProductID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "BuyProductID and GetProductID are required for buy_x_get_y promotion"})
			return
		}
	} else if existingPromotion.PromotionType == "percentage_discount" || existingPromotion.PromotionType == "fixed_discount" {
		if existingPromotion.DiscountValue <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "DiscountValue must be greater than 0 for discount promotions"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid promotion type"})
		return
	}

	if existingPromotion.EndDate.Before(existingPromotion.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "End date cannot be before start date"})
		return
	}

	if err := database.DB.Save(&existingPromotion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not update promotion", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Promotion updated", "data": existingPromotion})
}

// DeletePromotion handles deleting a promotion
func DeletePromotion(c *gin.Context) {
	id := c.Param("id")
	var promotion models.Promotion
	if err := database.DB.First(&promotion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Promotion not found", "data": err.Error()})
		return
	}

	if err := database.DB.Delete(&promotion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not delete promotion", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Promotion deleted"})
}
