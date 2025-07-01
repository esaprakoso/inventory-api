package handlers

import (
	"inventory/database"
	"inventory/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	var categories []models.Category
	database.DB.Find(&categories)
	c.JSON(http.StatusOK, categories)
}

func StoreCategory(c *gin.Context) {
	var data map[string]string

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	category := models.Category{
		Name: data["name"],
	}

	database.DB.Where("name = ?", category.Name).First(&category)
	if category.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Category already exists"})
		return
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

	var data map[string]string

	if err := c.BindJSON(&data); err != nil {
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

	category.Name = data["name"]

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
