package migrations

import (
	"fmt"
	"pos/database"
	"pos/models"
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
