package handlers

import (
	"inventory/database"
	"inventory/helpers"
	"inventory/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllProducts(c *gin.Context) {
	var products []models.Product
	database.DB.Preload("Category").Find(&products)
	c.JSON(http.StatusOK, products)
}

func StoreProduct(c *gin.Context) {
	var data map[string]any

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	price, err := helpers.GetFloat64(data["price"])
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": err,
		})
		return
	}

	category_id, err := helpers.GetInt(data["category_id"])
	if err != nil && data["category_id"] != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": err,
		})
		return
	}

	var categoryID *uint
	if category_id != 0 {
		categoryID = &[]uint{uint(category_id)}[0]
	}

	product := models.Product{
		Name:       data["name"].(string),
		SKU:        data["sku"].(string),
		Price:      price,
		CategoryID: categoryID,
	}

	database.DB.Where("SKU = ?", product.SKU).First(&product)
	if product.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Product SKU already exists"})
		return
	}

	database.DB.Preload("Category").Create(&product)

	c.JSON(http.StatusOK, product)
}

func GetProductByID(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	database.DB.Preload("Category").First(&product, id)

	if product.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

func UpdateProductByID(c *gin.Context) {
	id := c.Param("id")

	var data map[string]any

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var product models.Product
	database.DB.Preload("Category").First(&product, id)

	var productCek models.Product
	database.DB.First(&productCek, "SKU = ? and id != ?", data["sku"], id)

	if productCek.ID != 0 {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Product SKU already exists",
		})
		return
	}

	if product.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}

	price, err := helpers.GetFloat64(data["price"])
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": err,
		})
		return
	}

	category_id, err := helpers.GetInt(data["category_id"])
	if err != nil && data["category_id"] != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": err,
		})
		return
	}

	product.Name = data["name"].(string)
	product.SKU = data["sku"].(string)
	product.Price = price

	var updatedCategoryID *uint
	if category_id != 0 {
		updatedCategoryID = &[]uint{uint(category_id)}[0]
	}
	product.CategoryID = updatedCategoryID

	database.DB.Save(&product)

	c.JSON(http.StatusOK, product)
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
	c.JSON(http.StatusNotFound, gin.H{
		"message": "Product deleted",
	})
}
