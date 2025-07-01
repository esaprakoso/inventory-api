package middleware

import (
	"inventory/database"
	"inventory/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckUserExists() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data map[string]string

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request body",
			})
			c.Abort()
			return
		}

		var user models.User
		database.DB.Where("username = ?", data["username"]).First(&user)

		if user.ID != 0 {
			c.JSON(http.StatusConflict, gin.H{
				"message": "User with this username already exists",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
