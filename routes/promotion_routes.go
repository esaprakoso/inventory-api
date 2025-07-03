package routes

import (
	"pos/handlers"
	"pos/middleware"

	"github.com/gin-gonic/gin"
)

func SetupPromotionRoutes(router *gin.RouterGroup) {
	// Product Promotion routes
	router.POST("/product-promotions", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.CreateProductPromotion)
	router.GET("/product-promotions", middleware.Protected(), handlers.GetProductPromotions)
	router.GET("/product-promotions/:id", middleware.Protected(), handlers.GetProductPromotion)
	router.PUT("/product-promotions/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.UpdateProductPromotion)
	router.DELETE("/product-promotions/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.DeleteProductPromotion)

	// Cart Promotion routes
	router.POST("/cart-promotions", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.CreateCartPromotion)
	router.GET("/cart-promotions", middleware.Protected(), handlers.GetCartPromotions)
	router.GET("/cart-promotions/:id", middleware.Protected(), handlers.GetCartPromotion)
	router.PUT("/cart-promotions/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.UpdateCartPromotion)
	router.DELETE("/cart-promotions/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.DeleteCartPromotion)
}
