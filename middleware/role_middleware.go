package middleware

import (
	"inventory/database"
	"inventory/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizeRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
			c.Abort()
			return
		}

		var user models.User
		database.DB.First(&user, userID)

		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
			c.Abort()
			return
		}

		if user.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"message": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
