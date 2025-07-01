package handlers

import (
	"inventory/database"
	"inventory/models"
	"net/http"

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
	user.Role = data["role"]
	if data["password"] != "" {
		password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
		user.Password = string(password)
	}

	database.DB.Save(&user)

	c.JSON(http.StatusOK, user)
}
