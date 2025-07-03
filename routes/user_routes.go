package routes

import (
	"pos/handlers"
	"pos/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.RouterGroup) {
	router.GET("/users", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.GetAllUsers)
	router.GET("/users/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.GetUserByID)
	router.PATCH("/users/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.UpdateUserByID)
	router.DELETE("/users/:id", middleware.Protected(), middleware.AuthorizeRole("admin"), handlers.DeleteUserByID)
}
