package handlers

import (
	"fmt"
	"net/http"
	"time"

	"pos/database"
	"pos/models"

	"github.com/gin-gonic/gin"
)

type CartPromotionInput struct {
	PromotionType         string    `json:"promotion_type"` // e.g., "percentage_discount", "fixed_discount"
	DiscountValue         float64   `json:"discount_value"`
	MinimumPurchaseAmount float64   `json:"minimum_purchase_amount"`
	StartDate             time.Time `json:"start_date"`
	EndDate               time.Time `json:"end_date"`
}

type CartPromotionsResponse struct {
	Data []models.CartPromotion `json:"data"`
}

type CartPromotionResponse struct {
	Data models.CartPromotion `json:"data"`
}

// CreateCartPromotion handles the creation of a new cart promotion
// @Summary Create a new cart promotion
// @Description Create a new cart promotion. Admin only.
// @Tags Promotions
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   promotion body    CartPromotionInput true "Promotion data"
// @Success 201 {object} models.MessageResponse
// @Failure 400 {object} models.MessageResponse
// @Failure 401 {object} models.MessageResponse
// @Failure 500 {object} models.MessageResponse
// @Router /cart-promotions [post]
func CreateCartPromotion(c *gin.Context) {
	var cartPromotion models.CartPromotion
	if err := c.ShouldBindJSON(&cartPromotion); err != nil {
		c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Invalid request body"})
		return
	}

	// Validate cart promotion type and data
	if cartPromotion.PromotionType != "percentage_discount" && cartPromotion.PromotionType != "fixed_discount" {
		c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Invalid promotion type. Must be 'percentage_discount' or 'fixed_discount'"})
		return
	}

	if cartPromotion.DiscountValue <= 0 {
		c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "DiscountValue must be greater than 0"})
		return
	}

	if cartPromotion.MinimumPurchaseAmount <= 0 {
		c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "MinimumPurchaseAmount must be greater than 0"})
		return
	}

	if cartPromotion.EndDate.Before(cartPromotion.StartDate) {
		c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "End date cannot be before start date"})
		return
	}

	if err := database.DB.Create(&cartPromotion).Error; err != nil {
		fmt.Printf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, models.MessageResponse{Message: "Could not create cart promotion"})
		return
	}

	c.JSON(http.StatusCreated, models.MessageResponse{Message: "Cart Promotion created"})
}

// GetCartPromotions handles fetching all cart promotions
// @Summary Get all cart promotions
// @Description Get a list of all cart promotions.
// @Tags Promotions
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} CartPromotionsResponse
// @Router /cart-promotions [get]
func GetCartPromotions(c *gin.Context) {
	var cartPromotions []models.CartPromotion
	database.DB.Find(&cartPromotions)
	c.JSON(http.StatusOK, CartPromotionsResponse{Data: cartPromotions})
}

// GetCartPromotion handles fetching a single cart promotion by ID
// @Summary Get a cart promotion by ID
// @Description Get a single cart promotion by its ID.
// @Tags Promotions
// @Produce  json
// @Security BearerAuth
// @Param   id      path    int     true        "Promotion ID"
// @Success 200 {object} CartPromotionResponse
// @Failure 404 {object} models.MessageResponse
// @Router /cart-promotions/{id} [get]
func GetCartPromotion(c *gin.Context) {
	id := c.Param("id")
	var cartPromotion models.CartPromotion
	if err := database.DB.First(&cartPromotion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.MessageResponse{Message: "Cart Promotion not found"})
		return
	}
	c.JSON(http.StatusOK, CartPromotionResponse{Data: cartPromotion})
}

// UpdateCartPromotion handles updating an existing cart promotion
// @Summary Update a cart promotion by ID
// @Description Update a cart promotion's details by its ID. Admin only.
// @Tags Promotions
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   id      path    int     true        "Promotion ID"
// @Param   promotion body    CartPromotionInput true "Promotion data to update"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.MessageResponse
// @Failure 404 {object} models.MessageResponse
// @Failure 500 {object} models.MessageResponse
// @Router /cart-promotions/{id} [put]
func UpdateCartPromotion(c *gin.Context) {
	id := c.Param("id")
	var cartPromotion models.CartPromotion
	if err := c.ShouldBindJSON(&cartPromotion); err != nil {
		c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Invalid request body"})
		return
	}

	var existingCartPromotion models.CartPromotion
	if err := database.DB.First(&existingCartPromotion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.MessageResponse{Message: "Cart Promotion not found"})
		return
	}

	existingCartPromotion.PromotionType = cartPromotion.PromotionType
	existingCartPromotion.DiscountValue = cartPromotion.DiscountValue
	existingCartPromotion.MinimumPurchaseAmount = cartPromotion.MinimumPurchaseAmount
	existingCartPromotion.StartDate = cartPromotion.StartDate
	existingCartPromotion.EndDate = cartPromotion.EndDate

	// Validate cart promotion type and data
	if existingCartPromotion.PromotionType != "percentage_discount" && existingCartPromotion.PromotionType != "fixed_discount" {
		c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Invalid promotion type. Must be 'percentage_discount' or 'fixed_discount'"})
		return
	}

	if existingCartPromotion.DiscountValue <= 0 {
		c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "DiscountValue must be greater than 0"})
		return
	}

	if existingCartPromotion.MinimumPurchaseAmount <= 0 {
		c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "MinimumPurchaseAmount must be greater than 0"})
		return
	}

	if existingCartPromotion.EndDate.Before(existingCartPromotion.StartDate) {
		c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "End date cannot be before start date"})
		return
	}

	if err := database.DB.Save(&existingCartPromotion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.MessageResponse{Message: "Could not update cart promotion"})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: "Cart Promotion updated"})
}

// DeleteCartPromotion handles deleting a cart promotion
// @Summary Delete a cart promotion by ID
// @Description Delete a cart promotion by its ID. Admin only.
// @Tags Promotions
// @Produce  json
// @Security BearerAuth
// @Param   id      path    int     true        "Promotion ID"
// @Success 200 {object} models.MessageResponse
// @Failure 404 {object} models.MessageResponse
// @Failure 500 {object} models.MessageResponse
// @Router /cart-promotions/{id} [delete]
func DeleteCartPromotion(c *gin.Context) {
	id := c.Param("id")
	var cartPromotion models.CartPromotion
	if err := database.DB.First(&cartPromotion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.MessageResponse{Message: "Cart Promotion not found"})
		return
	}

	if err := database.DB.Delete(&cartPromotion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.MessageResponse{Message: "Could not delete cart promotion"})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: "Cart Promotion deleted"})
}
