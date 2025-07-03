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

		// Calculate discounted price for the product
		discountedPrice, activePromotion := utils.CalculateDiscountedPrice(product)
		itemPrice := product.Price
		itemDiscount := 0.0

		if activePromotion != nil {
			// Apply discount if active promotion is percentage or fixed discount
			if activePromotion.PromotionType == "percentage_discount" || activePromotion.PromotionType == "fixed_discount" {
				itemPrice = discountedPrice
				itemDiscount = product.Price - discountedPrice
			}
		}

		orderItem := models.OrderItem{
			OrderID:         order.ID,
			ProductID:       product.ID,
			Quantity:        itemInput.Quantity,
			Price:           product.Price, // Original price
			DiscountedPrice: itemPrice,     // Price after item-specific discount
			ItemDiscount:    itemDiscount,  // Discount amount for this item
			IsFreeItem:      false,
		}
		orderItems = append(orderItems, orderItem)

		grossTotal += product.Price * float64(itemInput.Quantity)
		itemDiscountTotal += itemDiscount * float64(itemInput.Quantity)

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
	order.TotalAmount = order.SubTotal - order.CartDiscount // Assuming CartDiscount is 0 for now

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not update order totals", "data": err.Error()})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Order created", "data": order})
}

// GetOrders handles fetching all orders
func GetOrders(c *gin.Context) {
	var orders []models.Order
	database.DB.Preload("OrderItems.Product").Find(&orders)
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Orders fetched", "data": orders})
}

// GetOrderByID handles fetching a single order by ID
func GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := database.DB.Preload("OrderItems.Product").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Order not found", "data": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Order fetched", "data": order})
}
