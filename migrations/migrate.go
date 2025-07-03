package migrations

import (
	"fmt"
	"pos/database"
	"pos/models"
)

func Migrate() {
	database.DB.AutoMigrate(
		&models.User{},
		&models.Order{},
		&models.OrderItem{},
		&models.Category{},
		&models.Product{},
		&models.StockTransaction{},
		&models.ProductPromotion{},
		&models.CartPromotion{},
	)
	fmt.Println("Database Migrated")
}
