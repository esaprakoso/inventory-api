package handlers

import (
	"net/http"

	"pos/database"
	"pos/models"

	"github.com/gin-gonic/gin"
)

// CreatePromotion handles the creation of a new promotion
// @Summary Create a new product promotion
// @Description Create a new product promotion. Admin only.
// @Tags Promotions
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   promotion body    models.ProductPromotion true "Promotion data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /product-promotions [post]
func CreateProductPromotion(c *gin.Context) {
	var promotion models.ProductPromotion
	if err := c.ShouldBindJSON(&promotion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body", "data": err.Error()})
		return
	}

	// Validate promotion type and data
	switch promotion.PromotionType {
	case "buy_x_get_y":
		if promotion.BuyProductID == nil || promotion.GetProductID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "BuyProductID and GetProductID are required for buy_x_get_y promotion"})
			return
		}
	case "percentage_discount", "fixed_discount":
		if promotion.DiscountValue <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "DiscountValue must be greater than 0 for discount promotions"})
			return
		}
	case "bundle_price":
		if promotion.RequiredQuantity == nil || promotion.PromoPrice == nil || *promotion.RequiredQuantity <= 0 || *promotion.PromoPrice <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "RequiredQuantity and PromoPrice must be greater than 0 for bundle_price promotion"})
			return
		}
	default:
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

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Product Promotion created", "data": promotion})
}

// GetPromotions handles fetching all promotions
// @Summary Get all product promotions
// @Description Get a list of all product promotions.
// @Tags Promotions
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /product-promotions [get]
func GetProductPromotions(c *gin.Context) {
	var promotions []models.ProductPromotion
	database.DB.Find(&promotions)
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Promotions fetched", "data": promotions})
}

// GetPromotion handles fetching a single promotion by ID
// @Summary Get a product promotion by ID
// @Description Get a single product promotion by its ID.
// @Tags Promotions
// @Produce  json
// @Security BearerAuth
// @Param   id      path    int     true        "Promotion ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /product-promotions/{id} [get]
func GetProductPromotion(c *gin.Context) {
	id := c.Param("id")
	var promotion models.ProductPromotion
	if err := database.DB.First(&promotion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Promotion not found", "data": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Promotion fetched", "data": promotion})
}

// UpdatePromotion handles updating an existing promotion
// @Summary Update a product promotion by ID
// @Description Update a product promotion's details by its ID. Admin only.
// @Tags Promotions
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   id      path    int     true        "Promotion ID"
// @Param   promotion body    models.ProductPromotion true "Promotion data to update"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /product-promotions/{id} [put]
func UpdateProductPromotion(c *gin.Context) {
	id := c.Param("id")
	var promotion models.ProductPromotion
	if err := c.ShouldBindJSON(&promotion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body", "data": err.Error()})
		return
	}

	var existingPromotion models.ProductPromotion
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
	existingPromotion.RequiredQuantity = promotion.RequiredQuantity
	existingPromotion.PromoPrice = promotion.PromoPrice
	existingPromotion.StartDate = promotion.StartDate
	existingPromotion.EndDate = promotion.EndDate

	// Validate promotion type and data
	switch existingPromotion.PromotionType {
	case "buy_x_get_y":
		if existingPromotion.BuyProductID == nil || existingPromotion.GetProductID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "BuyProductID and GetProductID are required for buy_x_get_y promotion"})
			return
		}
	case "percentage_discount", "fixed_discount":
		if existingPromotion.DiscountValue <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "DiscountValue must be greater than 0 for discount promotions"})
			return
		}
	case "bundle_price":
		if existingPromotion.RequiredQuantity == nil || existingPromotion.PromoPrice == nil || *existingPromotion.RequiredQuantity <= 0 || *existingPromotion.PromoPrice <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "RequiredQuantity and PromoPrice must be greater than 0 for bundle_price promotion"})
			return
		}
	default:
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
// @Summary Delete a product promotion by ID
// @Description Delete a product promotion by its ID. Admin only.
// @Tags Promotions
// @Produce  json
// @Security BearerAuth
// @Param   id      path    int     true        "Promotion ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /product-promotions/{id} [delete]
func DeleteProductPromotion(c *gin.Context) {
	id := c.Param("id")
	var promotion models.ProductPromotion
	if err := database.DB.First(&promotion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Promotion not found", "data": err.Error()})
		return
	}

	if err := database.DB.Delete(&promotion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not delete promotion", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Product Promotion deleted"})
}
