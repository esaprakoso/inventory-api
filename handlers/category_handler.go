package handlers

import (
	"inventory/database"
	"inventory/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"inventory/utils"
)

func GetCategories(c *gin.Context) {
	var categories []models.Category
	database.DB.Find(&categories)
	c.JSON(http.StatusOK, categories)
}

func StoreCategory(c *gin.Context) {
	type CreateCategoryInput struct {
		Name string `json:"name" binding:"required"`
	}
	var data CreateCategoryInput

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	isDup, err := utils.IsDuplicate[models.Category](database.DB, "name", data.Name, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	if isDup {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Category already exists"})
		return
	}

	category := models.Category{
		Name: data.Name,
	}

	database.DB.Create(&category)

	c.JSON(http.StatusOK, category)
}

func GetCategoryByID(c *gin.Context) {
	id := c.Param("id")

	var category models.Category
	database.DB.First(&category, id)

	if category.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Category not found",
		})
		return
	}

	c.JSON(http.StatusOK, category)
}

func UpdateCategoryByID(c *gin.Context) {
	id := c.Param("id")

	type UpdateCategoryInput struct {
		Name string `json:"name" binding:"required"`
	}
	var data UpdateCategoryInput

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var category models.Category
	database.DB.First(&category, id)

	if category.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Category not found",
		})
		return
	}

	isDup, err := utils.IsDuplicate[models.Category](database.DB, "name", data.Name, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	if isDup {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Category name already exists"})
		return
	}

	category.Name = data.Name

	database.DB.Save(&category)

	c.JSON(http.StatusOK, category)
}

func DeleteCategoryByID(c *gin.Context) {
	id := c.Param("id")

	var category models.Category
	database.DB.First(&category, id)

	if category.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Category not found",
		})
		return
	}

	database.DB.Delete(&category, id)
	c.JSON(http.StatusNotFound, gin.H{
		"message": "Category deleted",
	})
}
