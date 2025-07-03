package routes

import (
	"pos/handlers"
	"pos/middleware"

	"github.com/gin-gonic/gin"
)

func SetupProductRoutes(router *gin.RouterGroup) {
	router.GET("/products", middleware.Protected(), handlers.GetAllProducts)
	router.GET("/products/:id", middleware.Protected(), handlers.GetProductByID)
	router.POST("/products", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.StoreProduct)
	router.PUT("/products/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.UpdateProductByID)
	router.DELETE("/products/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.DeleteProductByID)
	router.PATCH("/products/:id/stock", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.UpdateProductStock)
}
