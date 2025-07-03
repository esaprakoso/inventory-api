package utils

import (
	"time"

	"pos/models"
)

// CalculateDiscountedPrice calculates the discounted price of a product based on active promotions.
func CalculateDiscountedPrice(product models.Product) (float64, *models.Promotion) {
	now := time.Now()
	var activePromotion *models.Promotion
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
