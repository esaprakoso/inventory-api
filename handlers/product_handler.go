package handlers

import (
	"net/http"
	"pos/database"
	"pos/models"

	"github.com/gin-gonic/gin"

	"pos/utils"
)

func GetAllProducts(c *gin.Context) {
	type ProductResponse struct {
		models.Product
		CategoryName string `json:"category_name"`
		Quantity     int    `json:"quantity"`
	}

	var products []ProductResponse
	var total int64

	page, _ := utils.GetInt(c.DefaultQuery("page", "1"))
	limit, _ := utils.GetInt(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	query := database.DB.Model(&models.Product{}).
		Select("products.*, categories.name as category_name, COALESCE(stocks.quantity, 0) as quantity").
		Joins("left join categories on categories.id = products.category_id").
		Joins("left join stocks on stocks.product_id = products.id")

	query.Count(&total)
	query.Limit(limit).Offset(offset).Find(&products)

	c.JSON(http.StatusOK, gin.H{
		"data":  products,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func StoreProduct(c *gin.Context) {
	type CreateProductInput struct {
		Name       string  `json:"name" binding:"required"`
		Price      float64 `json:"price" binding:"required"`
		SKU        string  `json:"sku" binding:"required"`
		CategoryID *uint   `json:"category_id"`
	}

	var data CreateProductInput

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	isDup, err := utils.IsDuplicate[models.Product](database.DB, "SKU", data.SKU, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	if isDup {
		c.JSON(406, gin.H{"message": "Product SKU already exists"})
		return
	}

	product := models.Product{
		Name:       data.Name,
		SKU:        data.SKU,
		Price:      data.Price,
		CategoryID: data.CategoryID,
	}

	database.DB.Create(&product)

	c.JSON(http.StatusOK, gin.H{"message": "Product Created"})
}

func GetProductByID(c *gin.Context) {
	id := c.Param("id")

	type ProductDetailResponse struct {
		models.Product
		CategoryName string `json:"category_name"`
		Quantity     int    `json:"quantity"`
	}

	var product ProductDetailResponse
	database.DB.Model(&models.Product{}).
		Select("products.*, categories.name as category_name, COALESCE(stocks.quantity, 0) as quantity").
		Joins("left join categories on categories.id = products.category_id").
		Joins("left join stocks on stocks.product_id = products.id").
		Where("products.id = ?", id).
		First(&product)

	if product.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": product,
	})
}

func UpdateProductByID(c *gin.Context) {
	id := c.Param("id")

	type UpdateProductInput struct {
		Name       string  `json:"name" binding:"required"`
		Price      float64 `json:"price" binding:"required"`
		SKU        string  `json:"sku" binding:"required"`
		CategoryID *uint   `json:"category_id"`
	}
	var data UpdateProductInput

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var product models.Product
	database.DB.Preload("Category").First(&product, id)

	isDup, err := utils.IsDuplicate[models.Product](database.DB, "SKU", data.SKU, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	if isDup {
		c.JSON(406, gin.H{"message": "Product SKU already exists"})
		return
	}

	if product.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}

	product.Name = data.Name
	product.SKU = data.SKU
	product.Price = data.Price
	product.CategoryID = data.CategoryID

	database.DB.Save(&product)

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated",
	})
}

func DeleteProductByID(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	database.DB.First(&product, id)

	if product.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}

	database.DB.Delete(&product, id)
	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted",
	})
}
