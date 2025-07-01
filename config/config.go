package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig(key string) string {
	godotenv.Load(".env")
	return os.Getenv(key)
}
