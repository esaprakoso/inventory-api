package migrations

import (
	"fmt"
	"inventory/database"
	"inventory/models"
)

func Migrate() {
	database.DB.AutoMigrate(
		&models.User{},
		
		&models.Category{},
		&models.Product{},
		&models.Stock{},
		&models.StockTransaction{},
	)
	fmt.Println("Database Migrated")
}
