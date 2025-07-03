package handlers

import (
	"fmt"
	"net/http"
	"time"

	"pos/database"
	"pos/models"
	"pos/utils"

	"github.com/gin-gonic/gin"
)

type CreateOrderInput struct {
	PaymentMethod string `json:"payment_method" binding:"required"`
	UserID        uint   `json:"user_id" binding:"required"`
	Items         []struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,min=1"`
	} `json:"items" binding:"required,min=1"`
}

// CreateOrder handles the creation of a new order
// @Summary Create a new order
// @Description Create a new order with specified products and payment method.
// @Tags Orders
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   order   body    CreateOrderInput true "Order details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders [post]
func CreateOrder(c *gin.Context) {
	var input CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body", "data": err.Error()})
		return
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to start transaction", "data": tx.Error.Error()})
		return
	}

	order := models.Order{
		PaymentMethod: input.PaymentMethod,
		UserID:        input.UserID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not create order", "data": err.Error()})
		return
	}

	var orderItems []models.OrderItem
	var grossTotal float64
	var itemDiscountTotal float64

	for _, itemInput := range input.Items {
		var product models.Product
		if err := tx.Preload("Promotions").First(&product, itemInput.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Product not found", "data": err.Error()})
			return
		}

		// Check stock availability
		if product.Quantity < itemInput.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Insufficient stock for product " + product.Name})
			return
		}

		// Decrement stock
		product.Quantity -= itemInput.Quantity
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to update stock", "data": err.Error()})
			return
		}

		// Create stock transaction log for 'out' type
		stockTransaction := models.StockTransaction{
			ProductID: product.ID,
			UserID:    order.UserID,
			Quantity:  itemInput.Quantity,
			Type:      models.StockTransactionTypeOut,
			SubType:   models.SubTypeSale,
			Notes:     fmt.Sprintf("Sale for order %d", order.ID),
		}
		if err := tx.Create(&stockTransaction).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create stock transaction", "data": err.Error()})
			return
		}

		// Calculate total price for the item, considering quantity and promotions
		totalItemPrice, activePromotion := utils.CalculateTotalPrice(product, itemInput.Quantity)
		originalItemTotal := product.Price * float64(itemInput.Quantity)
		itemDiscount := originalItemTotal - totalItemPrice

		orderItem := models.OrderItem{
			OrderID:         order.ID,
			ProductID:       product.ID,
			Quantity:        itemInput.Quantity,
			Price:           product.Price,  // Original price per unit
			DiscountedPrice: totalItemPrice, // Total price for all units of this item after discount
			ItemDiscount:    itemDiscount,   // Total discount for all units of this item
			IsFreeItem:      false,
		}
		orderItems = append(orderItems, orderItem)

		grossTotal += originalItemTotal
		itemDiscountTotal += itemDiscount

		// Handle "buy X get Y" promotion
		if activePromotion != nil && activePromotion.PromotionType == "buy_x_get_y" {
			if activePromotion.BuyProductID != nil && activePromotion.GetProductID != nil && *activePromotion.BuyProductID == product.ID {
				// Assuming for simplicity, 1 quantity of Y for 1 quantity of X
				var getProduct models.Product
				if err := tx.First(&getProduct, *activePromotion.GetProductID).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Get product for promotion not found", "data": err.Error()})
					return
				}
				// Add the 'get Y' product as a free item (Price and DiscountedPrice are 0)
				orderItems = append(orderItems, models.OrderItem{
					OrderID:         order.ID,
					ProductID:       getProduct.ID,
					Quantity:        itemInput.Quantity, // Same quantity as the 'buy X' product
					Price:           getProduct.Price,   // Original price of the free item
					DiscountedPrice: 0,                  // Free item
					ItemDiscount:    getProduct.Price,   // Discount is the full price of the item
					IsFreeItem:      true,
				})
			}
		}
	}

	if err := tx.Create(&orderItems).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not create order items", "data": err.Error()})
		return
	}

	order.GrossTotal = grossTotal
	order.ItemDiscountTotal = itemDiscountTotal
	order.SubTotal = grossTotal - itemDiscountTotal

	// Calculate cart-level discount
	order.CartDiscount = utils.CalculateCartDiscount(order.SubTotal)

	order.TotalAmount = order.SubTotal - order.CartDiscount

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not update order totals", "data": err.Error()})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Order created", "data": order})
}

// GetOrders handles fetching all orders
// @Summary Get all orders
// @Description Get a list of all orders. Admin can see all, users see their own.
// @Tags Orders
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /orders [get]
func GetOrders(c *gin.Context) {
	var orders []models.Order
	database.DB.Preload("OrderItems.Product").Find(&orders)
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Orders fetched", "data": orders})
}

// GetOrderByID handles fetching a single order by ID
// @Summary Get an order by ID
// @Description Get a single order by its ID. Admin can see any order, users see their own.
// @Tags Orders
// @Produce  json
// @Security BearerAuth
// @Param   id      path    int     true        "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id} [get]
func GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := database.DB.Preload("OrderItems.Product").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Order not found", "data": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Order fetched", "data": order})
}
