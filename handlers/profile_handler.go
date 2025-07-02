package handlers

import (
	"inventory/database"
	"inventory/models"
	"net/http"

	"inventory/utils"

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

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func UpdateProfile(c *gin.Context) {
	id := c.GetString("user_id")
	type UpdateProfileInput struct {
		Username string `json:"username" binding:"required"`
		Name     string `json:"name" binding:"required"`
	}
	var data UpdateProfileInput

	if err := c.ShouldBindJSON(&data); err != nil {
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

	isDup, err := utils.IsDuplicate[models.User](database.DB, "username", data.Username, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	if isDup {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Username already exists"})
		return
	}

	user.Username = data.Username
	user.Name = data.Name

	database.DB.Save(&user)

	c.JSON(http.StatusOK, user)
}

func UpdateProfilePassword(c *gin.Context) {
	id := c.GetString("user_id")
	type UpdatePasswordInput struct {
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	}
	var data UpdatePasswordInput

	if err := c.ShouldBindJSON(&data); err != nil {
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

	password, err := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
		return
	}
	user.Password = string(password)
	database.DB.Save(&user)
	c.JSON(http.StatusOK, user)

}
