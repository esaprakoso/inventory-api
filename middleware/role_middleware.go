package middleware

import (
	"net/http"
	"pos/database"
	"pos/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AuthorizeRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetString("user_id")
		if userIDStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
			c.Abort()
			return
		}

		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid User ID format in context"})
			c.Abort()
			return
		}

		var user models.User
		database.DB.First(&user, uint(userID))

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
