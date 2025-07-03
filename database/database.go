package database

import (
	"fmt"
	"log"

	"pos/config"

	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn, // ⬅️ Hanya log warning/error
			IgnoreRecordNotFoundError: true,        // ⬅️ Abaikan log untuk "record not found"
			Colorful:                  true,
		},
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	fmt.Println("Database connection successfully opened")

}
