package handlers

import (
	"inventory/database"
	"inventory/models"
	"net/http"

	"inventory/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers(c *gin.Context) {
	var users []models.User
	database.DB.Find(&users)
	c.JSON(http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	database.DB.First(&user, id)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUserByID(c *gin.Context) {
	id := c.Param("id")
	type UpdateUserInput struct {
		Username string  `json:"username" binding:"required"`
		Name     string  `json:"name" binding:"required"`
		Role     string  `json:"role" binding:"required"`
		Password *string `json:"password"`
	}

	var data UpdateUserInput

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user models.User
	database.DB.First(&user, id)

	isDup, _ := utils.CheckDuplicate[models.User](database.DB, "username", data.Username, id)
	if isDup {
		c.JSON(406, gin.H{"message": "Username already exists"})
		return
	}

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	user.Username = data.Username
	user.Name = data.Name
	user.Role = data.Role
	if data.Password != nil {
		password, _ := bcrypt.GenerateFromPassword([]byte(*data.Password), 14)
		user.Password = string(password)
	}

	database.DB.Save(&user)

	c.JSON(http.StatusOK, user)
}

func DeleteUserByID(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	database.DB.First(&user, id)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	database.DB.Delete(&user, id)
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted",
	})
}
