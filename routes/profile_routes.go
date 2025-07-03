package routes

import (
	"pos/handlers"
	"pos/middleware"

	"github.com/gin-gonic/gin"
)

func SetupProfileRoutes(router *gin.RouterGroup) {
	router.GET("/profile", middleware.Protected(), handlers.GetUserProfile)
	router.PATCH("/profile", middleware.Protected(), handlers.GetUserProfile)
	router.PATCH("/profile/password", middleware.Protected(), handlers.UpdateProfilePassword)
}
