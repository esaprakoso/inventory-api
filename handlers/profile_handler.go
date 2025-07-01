package handlers

import (
	"inventory/database"
	"inventory/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUserProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	var user models.User
	database.DB.First(&user, userID)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateProfile(c *gin.Context) {
	id := c.GetString("user_id")
	var data map[string]string

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user models.User
	database.DB.First(&user, id)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	user.Username = data["username"]
	user.Name = data["name"]

	database.DB.Save(&user)

	c.JSON(http.StatusOK, user)
}

func UpdateProfilePassword(c *gin.Context) {
	id := c.GetString("user_id")
	var data map[string]string

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user models.User
	database.DB.First(&user, id)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	if data["password"] != data["confirm_password"] {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Password not match",
		})
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	user.Password = string(password)
	database.DB.Save(&user)
	c.JSON(http.StatusOK, user)

}
