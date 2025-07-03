package handlers

import (
	"net/http"
	"pos/database"
	"pos/models"

	"pos/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UpdateUserInput struct {
	Username string  `json:"username" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Role     string  `json:"role" binding:"required,oneof=admin user"`
	Password *string `json:"password"`
}

// @Summary Get all users
// @Description Get a list of all users. Admin only.
// @Tags Users
// @Produce  json
// @Security BearerAuth
// @Param   page      query    int     false        "Page number"
// @Param   limit     query    int     false        "Number of items per page"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users [get]
func GetAllUsers(c *gin.Context) {
	var users []models.User
	var total int64

	page, _ := utils.GetInt(c.DefaultQuery("page", "1"))
	limit, _ := utils.GetInt(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	database.DB.Model(&models.User{}).Count(&total)
	database.DB.Limit(limit).Offset(offset).Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"data":  users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// @Summary Get a user by ID
// @Description Get a single user by their ID. Admin only.
// @Tags Users
// @Produce  json
// @Security BearerAuth
// @Param   id      path    int     true        "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id} [get]
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

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// @Summary Update a user by ID
// @Description Update a user's details by their ID. Admin only.
// @Tags Users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   id      path    int     true        "User ID"
// @Param   user    body    UpdateUserInput true    "User data to update"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 406 {object} map[string]interface{}
// @Router /users/{id} [patch]
func UpdateUserByID(c *gin.Context) {
	id := c.Param("id")
	var data UpdateUserInput

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user models.User
	database.DB.First(&user, id)

	isDup, err := utils.IsDuplicate[models.User](database.DB, "username", data.Username, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
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
		password, err := bcrypt.GenerateFromPassword([]byte(*data.Password), 14)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
			return
		}
		user.Password = string(password)
	}

	database.DB.Save(&user)

	c.JSON(http.StatusOK, user)
}

// @Summary Delete a user by ID
// @Description Delete a user by their ID. Admin only.
// @Tags Users
// @Produce  json
// @Security BearerAuth
// @Param   id      path    int     true        "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id} [delete]
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
