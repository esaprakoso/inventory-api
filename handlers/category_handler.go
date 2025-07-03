package handlers

import (
	"net/http"
	"pos/database"
	"pos/models"

	"pos/utils"

	"github.com/gin-gonic/gin"
)

type CategoryInput struct {
	Name string `json:"name" binding:"required"`
}

type CategoriesResponse struct {
	Data  []models.Category `json:"data"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
}

type CategoryResponse struct {
	Data models.Category `json:"data"`
}

// @Summary Get all categories
// @Description Get a list of all categories.
// @Tags Categories
// @Produce  json
// @Security BearerAuth
// @Param   page      query    int     false        "Page number"
// @Param   limit     query    int     false        "Number of items per page"
// @Success 200 {object} CategoriesResponse
// @Router /categories [get]
func GetCategories(c *gin.Context) {
	var categories []models.Category
	var total int64

	page, _ := utils.GetInt(c.DefaultQuery("page", "1"))
	limit, _ := utils.GetInt(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	database.DB.Model(&models.Category{}).Count(&total)
	database.DB.Limit(limit).Offset(offset).Find(&categories)

	c.JSON(http.StatusOK, CategoriesResponse{
		Data:  categories,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

// @Summary Create a new category
// @Description Create a new category. Admin only.
// @Tags Categories
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   category body    CategoryInput true "Category data"
// @Success 201 {object} models.MessageResponse
// @Failure 400 {object} models.MessageResponse
// @Failure 401 {object} models.MessageResponse
// @Failure 406 {object} models.MessageResponse
// @Failure 500 {object} models.MessageResponse
// @Router /categories [post]
func StoreCategory(c *gin.Context) {
	var data CategoryInput

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Invalid request body"})
		return
	}

	isDup, err := utils.IsDuplicate[models.Category](database.DB, "name", data.Name, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.MessageResponse{Message: "Database error"})
		return
	}
	if isDup {
		c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Category already exists"})
		return
	}

	category := models.Category{
		Name: data.Name,
	}

	database.DB.Create(&category)

	c.JSON(http.StatusOK, models.MessageResponse{Message: "Category created"})
}

// @Summary Get a category by ID
// @Description Get a single category by its ID.
// @Tags Categories
// @Produce  json
// @Security BearerAuth
// @Param   id      path    int     true        "Category ID"
// @Success 200 {object} CategoryResponse
// @Failure 404 {object} models.MessageResponse
// @Router /categories/{id} [get]
func GetCategoryByID(c *gin.Context) {
	id := c.Param("id")

	var category models.Category
	database.DB.First(&category, id)

	if category.ID == 0 {
		c.JSON(http.StatusNotFound, models.MessageResponse{
			Message: "Category not found",
		})
		return
	}

	c.JSON(http.StatusOK, CategoryResponse{
		Data: category,
	})
}

// @Summary Update a category by ID
// @Description Update a category's details by its ID. Admin only.
// @Tags Categories
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   id      path    int     true        "Category ID"
// @Param   category body    CategoryInput true "Category data to update"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.MessageResponse
// @Failure 401 {object} models.MessageResponse
// @Failure 404 {object} models.MessageResponse
// @Failure 406 {object} models.MessageResponse
// @Router /categories/{id} [put]
func UpdateCategoryByID(c *gin.Context) {
	id := c.Param("id")
	var data CategoryInput

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, models.MessageResponse{Message: err.Error()})
		return
	}

	var category models.Category
	database.DB.First(&category, id)

	if category.ID == 0 {
		c.JSON(http.StatusNotFound, models.MessageResponse{
			Message: "Category not found",
		})
		return
	}

	isDup, err := utils.IsDuplicate[models.Category](database.DB, "name", data.Name, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.MessageResponse{Message: "Database error"})
		return
	}
	if isDup {
		c.JSON(http.StatusBadRequest, models.MessageResponse{Message: "Category name already exists"})
		return
	}

	category.Name = data.Name

	database.DB.Save(&category)

	c.JSON(http.StatusOK, models.MessageResponse{Message: "Category updated"})
}

// @Summary Delete a category by ID
// @Description Delete a category by its ID. Admin only.
// @Tags Categories
// @Produce  json
// @Security BearerAuth
// @Param   id      path    int     true        "Category ID"
// @Success 200 {object} models.MessageResponse
// @Failure 401 {object} models.MessageResponse
// @Failure 404 {object} models.MessageResponse
// @Router /categories/{id} [delete]
func DeleteCategoryByID(c *gin.Context) {
	id := c.Param("id")

	var category models.Category
	database.DB.First(&category, id)

	if category.ID == 0 {
		c.JSON(http.StatusNotFound, models.MessageResponse{
			Message: "Category not found",
		})
		return
	}

	database.DB.Delete(&category, id)
	c.JSON(http.StatusNotFound, models.MessageResponse{
		Message: "Category deleted",
	})
}
