package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(app *gin.Engine) {
	api := app.Group("/api")

	SetupAuthRoutes(api)
	SetupUserRoutes(api)
	SetupProfileRoutes(api)
	SetupProductRoutes(api)
	SetupCategoryRoutes(api)
	SetupPromotionRoutes(api)
	SetupOrderRoutes(api)
}