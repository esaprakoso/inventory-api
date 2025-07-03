package handlers

import (
	"net/http"
	"time"

	"pos/database"
	"pos/models"

	"github.com/gin-gonic/gin"
)

type ProductPromotionInput struct {
	ProductID        uint      `json:"product_id"`
	PromotionType    string    `json:"promotion_type"` // e.g., "percentage_discount", "fixed_discount", "buy_x_get_y", "bundle_price"
	DiscountValue    float64   `json:"discount_value,omitempty"`
	BuyProductID     *uint     `json:"buy_product_id,omitempty"`
	GetProductID     *uint     `json:"get_product_id,omitempty"`
	RequiredQuantity *int      `json:"required_quantity,omitempty"` // For "bundle_price" or "buy_x_get_y"
	PromoPrice       *float64  `json:"promo_price,omitempty"`       // For "bundle_price"
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
}

type ProductPromotionsResponse struct {
	Data []models.ProductPromotion `json:"data"`
}

type ProductPromotionResponse struct {
	Data models.ProductPromotion `json:"data"`
}

// CreatePromotion handles the creation of a new promotion
// @Summary Create a new product promotion
// @Description Create a new product promotion. Admin only.
// @Tags Promotions
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   promotion body    ProductPromotionInput true "Promotion data"
// @Success 201 {object} models.MessageResponse
// @Failure 400 {object} models.MessageResponse
// @Failure 401 {object} models.MessageResponse
// @Failure 500 {object} models.MessageResponse
// @Router /product-promotions [post]
func CreateProductPromotion(c *gin.Context) {
	var promotion ProductPromotionInput
	if err := c.ShouldBindJSON(&promotion); err != nil {
		c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Invalid request body"})
		return
	}

	// Validate promotion type and data
	switch promotion.PromotionType {
	case "buy_x_get_y":
		if promotion.BuyProductID == nil || promotion.GetProductID == nil {
			c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "BuyProductID and GetProductID are required for buy_x_get_y promotion"})
			return
		}
	case "percentage_discount", "fixed_discount":
		if promotion.DiscountValue <= 0 {
			c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "DiscountValue must be greater than 0 for discount promotions"})
			return
		}
	case "bundle_price":
		if promotion.RequiredQuantity == nil || promotion.PromoPrice == nil || *promotion.RequiredQuantity <= 0 || *promotion.PromoPrice <= 0 {
			c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "RequiredQuantity and PromoPrice must be greater than 0 for bundle_price promotion"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Invalid promotion type"})
		return
	}

	if promotion.EndDate.Before(promotion.StartDate) {
		c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "End date cannot be before start date"})
		return
	}

	if err := database.DB.Create(&promotion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.MessageResponse{Message: "Could not create promotion"})
		return
	}

	c.JSON(http.StatusCreated, models.MessageResponse{Message: "Product Promotion created"})
}

// GetPromotions handles fetching all promotions
// @Summary Get all product promotions
// @Description Get a list of all product promotions.
// @Tags Promotions
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} ProductPromotionsResponse
// @Router /product-promotions [get]
func GetProductPromotions(c *gin.Context) {
	var promotions []models.ProductPromotion
	database.DB.Find(&promotions)
	c.JSON(http.StatusOK, ProductPromotionsResponse{Data: promotions})
}

// GetPromotion handles fetching a single promotion by ID
// @Summary Get a product promotion by ID
// @Description Get a single product promotion by its ID.
// @Tags Promotions
// @Produce  json
// @Security BearerAuth
// @Param   id      path    int     true        "Promotion ID"
// @Success 200 {object} ProductPromotionResponse
// @Failure 404 {object} models.MessageResponse
// @Router /product-promotions/{id} [get]
func GetProductPromotion(c *gin.Context) {
	id := c.Param("id")
	var promotion models.ProductPromotion
	if err := database.DB.First(&promotion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.MessageResponse{Message: "Promotion not found"})
		return
	}
	c.JSON(http.StatusOK, ProductPromotionResponse{Data: promotion})
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
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.MessageResponse
// @Failure 404 {object} models.MessageResponse
// @Router /product-promotions/{id} [put]
func UpdateProductPromotion(c *gin.Context) {
	id := c.Param("id")
	var promotion models.ProductPromotion
	if err := c.ShouldBindJSON(&promotion); err != nil {
		c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Invalid request body"})
		return
	}

	var existingPromotion models.ProductPromotion
	if err := database.DB.First(&existingPromotion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.MessageResponse{Message: "Promotion not found"})
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
			c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "BuyProductID and GetProductID are required for buy_x_get_y promotion"})
			return
		}
	case "percentage_discount", "fixed_discount":
		if existingPromotion.DiscountValue <= 0 {
			c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "DiscountValue must be greater than 0 for discount promotions"})
			return
		}
	case "bundle_price":
		if existingPromotion.RequiredQuantity == nil || existingPromotion.PromoPrice == nil || *existingPromotion.RequiredQuantity <= 0 || *existingPromotion.PromoPrice <= 0 {
			c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "RequiredQuantity and PromoPrice must be greater than 0 for bundle_price promotion"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Invalid promotion type"})
		return
	}

	if existingPromotion.EndDate.Before(existingPromotion.StartDate) {
		c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "End date cannot be before start date"})
		return
	}

	if err := database.DB.Save(&existingPromotion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.MessageResponse{Message: "Could not update promotion"})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: "Promotion updated"})
}

// DeletePromotion handles deleting a promotion
// @Summary Delete a product promotion by ID
// @Description Delete a product promotion by its ID. Admin only.
// @Tags Promotions
// @Produce  json
// @Security BearerAuth
// @Param   id      path    int     true        "Promotion ID"
// @Success 200 {object} models.MessageResponse
// @Failure 401 {object} models.MessageResponse
// @Failure 404 {object} models.MessageResponse
// @Router /product-promotions/{id} [delete]
func DeleteProductPromotion(c *gin.Context) {
	id := c.Param("id")
	var promotion models.ProductPromotion
	if err := database.DB.First(&promotion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.MessageResponse{Message: "Promotion not found"})
		return
	}

	if err := database.DB.Delete(&promotion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.MessageResponse{Message: "Could not delete promotion"})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: "Product Promotion deleted"})
}
