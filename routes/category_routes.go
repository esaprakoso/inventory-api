package routes

import (
	"pos/handlers"
	"pos/middleware"

	"github.com/gin-gonic/gin"
)

func SetupCategoryRoutes(router *gin.RouterGroup) {
	router.GET("/categories", middleware.Protected(), handlers.GetCategories)
	router.GET("/categories/:id", middleware.Protected(), handlers.GetCategoryByID)
	router.POST("/categories", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.StoreCategory)
	router.PUT("/categories/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.UpdateCategoryByID)
	router.DELETE("/categories/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.DeleteCategoryByID)
}
