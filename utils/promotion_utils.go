package utils

import (
	"time"

	"pos/database"
	"pos/models"
)

// CalculateDiscountedPrice calculates the discounted price of a product based on active promotions.
func CalculateDiscountedPrice(product models.Product) (float64, *models.ProductPromotion) {
	now := time.Now()
	var activePromotion *models.ProductPromotion
	discountedPrice := product.Price

	for i := range product.Promotions {
		p := &product.Promotions[i]
		if now.After(p.StartDate) && now.Before(p.EndDate) {
			activePromotion = p
			break
		}
	}

	if activePromotion != nil {
		switch activePromotion.PromotionType {
		case "percentage_discount":
			discountedPrice = product.Price * (1 - activePromotion.DiscountValue/100)
		case "fixed_discount":
			discountedPrice = product.Price - activePromotion.DiscountValue
		}
		if discountedPrice < 0 {
			discountedPrice = 0
		}
	}

	return discountedPrice, activePromotion
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
