package routes

import (
	"pos/handlers"
	"pos/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(app *gin.Engine) {
	api := app.Group("/api")

	// Product Promotion routes
	api.POST("/product-promotions", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.CreateProductPromotion)
	api.GET("/product-promotions", middleware.Protected(), handlers.GetProductPromotions)
	api.GET("/product-promotions/:id", middleware.Protected(), handlers.GetProductPromotion)
	api.PUT("/product-promotions/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.UpdateProductPromotion)
	api.DELETE("/product-promotions/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.DeleteProductPromotion)

	// Cart Promotion routes
	api.POST("/cart-promotions", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.CreateCartPromotion)
	api.GET("/cart-promotions", middleware.Protected(), handlers.GetCartPromotions)
	api.GET("/cart-promotions/:id", middleware.Protected(), handlers.GetCartPromotion)
	api.PUT("/cart-promotions/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.UpdateCartPromotion)
	api.DELETE("/cart-promotions/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.DeleteCartPromotion)

	// Order routes
	api.POST("/orders", middleware.Protected(), handlers.CreateOrder)
	api.GET("/orders", middleware.Protected(), handlers.GetOrders)
	api.GET("/orders/:id", middleware.Protected(), handlers.GetOrderByID)

	// Auth routes
	auth := api.Group("/auth")
	auth.POST("/register", handlers.Register)
	auth.POST("/login", handlers.Login)
	auth.POST("/refresh", handlers.RefreshToken)

	// User routes
	api.GET("/users", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.GetAllUsers)
	api.GET("/users/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.GetUserByID)
	api.PATCH("/users/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.UpdateUserByID)
	api.DELETE("/users/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.DeleteUserByID)

	// Profile routes
	api.GET("/profile", middleware.Protected(), handlers.GetUserProfile)
	api.PATCH("/profile", middleware.Protected(), handlers.GetUserProfile)
	api.PATCH("/profile/password", middleware.Protected(), handlers.UpdateProfilePassword)

	// Product routes
	api.GET("/products", middleware.Protected(), handlers.GetAllProducts)
	api.GET("/products/:id", middleware.Protected(), handlers.GetProductByID)
	api.POST("/products", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.StoreProduct)
	api.PUT("/products/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.UpdateProductByID)
	api.DELETE("/products/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.DeleteProductByID)
	api.PATCH("/products/:id/stock", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.UpdateProductStock)

	// Category routes
	api.GET("/categories", middleware.Protected(), handlers.GetCategories)
	api.GET("/categories/:id", middleware.Protected(), handlers.GetCategoryByID)
	api.POST("/categories", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.StoreCategory)
	api.PUT("/categories/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.UpdateCategoryByID)
	api.DELETE("/categories/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.DeleteCategoryByID)
}
