package main

import (
	"inventory/database"
	"inventory/migrations"
)

func main() {
	database.Connect()
	migrations.Migrate()
}
