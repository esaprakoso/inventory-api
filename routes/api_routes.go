package routes

import (
	"inventory/handlers"
	"inventory/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(app *gin.Engine) {
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.POST("/register", middleware.CheckUserExists(), handlers.Register)
	auth.POST("/login", handlers.Login)
	auth.POST("/refresh", handlers.RefreshToken)

	// User routes
	api.GET("/users", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.GetAllUsers)
	api.GET("/users/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.GetUserByID)
	api.PATCH("/users/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.UpdateUserByID)

	// Profile routes
	api.GET("/profile", middleware.Protected(), handlers.GetUserProfile)
	api.PATCH("/profile", middleware.Protected(), handlers.GetUserProfile)
	api.PATCH("/profile/password", middleware.Protected(), handlers.UpdateProfilePassword)

	// Warehouse routes
	api.GET("/warehouse", middleware.Protected(), handlers.GetAllWarehouses)
	api.GET("/warehouse/:id", middleware.Protected(), handlers.GetWarehouseByID)
	api.POST("/warehouse", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.StoreWarehouse)
	api.PUT("/warehouse/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.UpdateWarehouseByID)
	api.DELETE("/warehouse/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.DeleteWarehouseByID)

	// Product routes
	api.GET("/products", middleware.Protected(), handlers.GetAllProducts)
	api.GET("/products/:id", middleware.Protected(), handlers.GetProductByID)
	api.POST("/products", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.StoreProduct)
	api.PUT("/products/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.UpdateProductByID)
	api.DELETE("/products/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.DeleteProductByID)

	// Product routes
	api.GET("/categories", middleware.Protected(), handlers.GetCategories)
	api.GET("/categories/:id", middleware.Protected(), handlers.GetCategoryByID)
	api.POST("/categories", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.StoreCategory)
	api.PUT("/categories/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.UpdateCategoryByID)
	api.DELETE("/categories/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.DeleteCategoryByID)
}
