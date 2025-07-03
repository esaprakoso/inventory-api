package utils

import (
	"time"

	"pos/database"
	"pos/models"
)

// CalculateTotalPrice calculates the total price for a given quantity of a product, applying the best active promotion.
func CalculateTotalPrice(product models.Product, quantity int) (float64, *models.ProductPromotion) {
	now := time.Now()
	var activePromotion *models.ProductPromotion

	// Find the first active promotion for the product
	for i := range product.Promotions {
		p := &product.Promotions[i]
		if now.After(p.StartDate) && now.Before(p.EndDate) {
			activePromotion = p
			break
		}
	}

	totalPrice := product.Price * float64(quantity)

	if activePromotion != nil {
		switch activePromotion.PromotionType {
		case "bundle_price":
			if activePromotion.RequiredQuantity != nil && activePromotion.PromoPrice != nil && *activePromotion.RequiredQuantity > 0 {
				numBundles := quantity / *activePromotion.RequiredQuantity
				remainingItems := quantity % *activePromotion.RequiredQuantity
				totalPrice = float64(numBundles)*(*activePromotion.PromoPrice) + float64(remainingItems)*product.Price
			}
		case "percentage_discount":
			discountedPrice := product.Price * (1 - activePromotion.DiscountValue/100)
			totalPrice = discountedPrice * float64(quantity)
		case "fixed_discount":
			discountedPrice := product.Price - activePromotion.DiscountValue
			totalPrice = discountedPrice * float64(quantity)
		}
		if totalPrice < 0 {
			totalPrice = 0
		}
	}

	return totalPrice, activePromotion
}

// CalculateCartDiscount calculates a cart-level discount based on total purchase amount.
func CalculateCartDiscount(subTotal float64) float64 {
	var activeCartPromotion models.CartPromotion
	now := time.Now()

	// Find an active cart promotion that applies to the current subTotal
	database.DB.Where(
		"minimum_purchase_amount <= ? AND start_date <= ? AND end_date >= ?",
		subTotal, now, now,
	).First(&activeCartPromotion)

	if activeCartPromotion.ID != 0 {
		switch activeCartPromotion.PromotionType {
		case "percentage_discount":
			return subTotal * (activeCartPromotion.DiscountValue / 100)
		case "fixed_discount":
			return activeCartPromotion.DiscountValue
		}
	}
	return 0.0
}
