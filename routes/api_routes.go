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
	api.GET("/warehouse", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.GetAllWarehouses)
	api.POST("/warehouse", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.StoreWarehouse)
	api.GET("/warehouse/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.GetWarehouseByID)
	api.PUT("/warehouse/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.UpdateWarehouseByID)
}
