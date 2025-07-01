package handlers

import (
	"inventory/database"
	"inventory/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllWarehouses(c *gin.Context) {
	OwnerID := c.GetString("user_id")
	var warehouses []models.Warehouse
	database.DB.Where("owner_id = ?", OwnerID).Find(&warehouses)
	c.JSON(http.StatusOK, warehouses)
}

func StoreWarehouse(c *gin.Context) {
	OwnerID := c.GetString("user_id")
	var data map[string]string

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	warehouse := models.Warehouse{
		Name:     data["name"],
		Location: data["location"],
		OwnerID: func() uint {
			id, err := strconv.ParseUint(OwnerID, 10, 64)
			if err != nil {
				return 0 // Or handle the error appropriately, e.g., return an error response
			}
			return uint(id)
		}(),
	}

	database.DB.Create(&warehouse)

	c.JSON(http.StatusOK, warehouse)
}

func GetWarehouseByID(c *gin.Context) {
	id := c.Param("id")

	var warehouse models.Warehouse
	database.DB.First(&warehouse, id)

	if warehouse.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Warehouse not found",
		})
		return
	}

	c.JSON(http.StatusOK, warehouse)
}

func UpdateWarehouseByID(c *gin.Context) {
	id := c.Param("id")

	var data map[string]string

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var warehouse models.Warehouse
	database.DB.First(&warehouse, id)

	if warehouse.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Warehouse not found",
		})
		return
	}

	warehouse.Name = data["name"]
	warehouse.Location = data["location"]

	database.DB.Save(&warehouse)

	c.JSON(http.StatusOK, warehouse)
}

func DeleteWarehouseByID(c *gin.Context) {
	id := c.Param("id")

	var warehouse models.Warehouse
	database.DB.First(&warehouse, id)

	if warehouse.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Warehouse not found",
		})
		return
	}

	database.DB.Delete(&warehouse, id)
	c.JSON(http.StatusNotFound, gin.H{
		"message": "Warehouse deleted",
	})
}
