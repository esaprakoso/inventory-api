package handlers

import (
	"net/http"
	"pos/database"
	"pos/models"

	"pos/utils"

	"github.com/gin-gonic/gin"
)

type CreateCategoryInput struct {
	Name string `json:"name" binding:"required"`
}

type UpdateCategoryInput struct {
	Name string `json:"name" binding:"required"`
}

// @Summary Get all categories
// @Description Get a list of all categories.
// @Tags Categories
// @Produce  json
// @Security BearerAuth
// @Param   page      query    int     false        "Page number"
// @Param   limit     query    int     false        "Number of items per page"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /categories [get]
func GetCategories(c *gin.Context) {
	var categories []models.Category
	var total int64

	page, _ := utils.GetInt(c.DefaultQuery("page", "1"))
	limit, _ := utils.GetInt(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	database.DB.Model(&models.Category{}).Count(&total)
	database.DB.Limit(limit).Offset(offset).Find(&categories)

	c.JSON(http.StatusOK, gin.H{
		"data":  categories,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// @Summary Create a new category
// @Description Create a new category. Admin only.
// @Tags Categories
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   category body    CreateCategoryInput true "Category data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 406 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /categories [post]
func StoreCategory(c *gin.Context) {
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

// @Summary Get a category by ID
// @Description Get a single category by its ID.
// @Tags Categories
// @Produce  json
// @Security BearerAuth
// @Param   id      path    int     true        "Category ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /categories/{id} [get]
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

	c.JSON(http.StatusOK, gin.H{
		"data": category,
	})
}

// @Summary Update a category by ID
// @Description Update a category's details by its ID. Admin only.
// @Tags Categories
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   id      path    int     true        "Category ID"
// @Param   category body    UpdateCategoryInput true "Category data to update"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 406 {object} map[string]interface{}
// @Router /categories/{id} [put]
func UpdateCategoryByID(c *gin.Context) {
	id := c.Param("id")
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

// @Summary Delete a category by ID
// @Description Delete a category by its ID. Admin only.
// @Tags Categories
// @Produce  json
// @Security BearerAuth
// @Param   id      path    int     true        "Category ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /categories/{id} [delete]
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
