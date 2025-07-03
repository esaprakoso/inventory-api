package main

import (
	"pos/database"
	"pos/migrations"
)

func main() {
	database.Connect()
	migrations.Migrate()
}
