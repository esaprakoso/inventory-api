package database

import (
	"fmt"
	"log"

	"inventory/config"
	"inventory/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.LoadConfig("DB_HOST"),
		config.LoadConfig("DB_USER"),
		config.LoadConfig("DB_PASSWORD"),
		config.LoadConfig("DB_NAME"),
		config.LoadConfig("DB_PORT"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	fmt.Println("Database connection successfully opened")

	DB.AutoMigrate(
		&models.User{},
		&models.Warehouse{},
	)
	fmt.Println("Database Migrated")
}
