package routes

import (
	"pos/handlers"
	"pos/middleware"

	"github.com/gin-gonic/gin"
)

func SetupOrderRoutes(router *gin.RouterGroup) {
	router.POST("/orders", middleware.Protected(), handlers.CreateOrder)
	router.GET("/orders", middleware.Protected(), handlers.GetOrders)
	router.GET("/orders/:id", middleware.Protected(), handlers.GetOrderByID)
}
