package main

import (
	"inventory/database"
	"inventory/routes"
	"log"

	"inventory/validators"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {
	app := gin.Default()

	database.Connect()

	// daftar validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validators.RegisterCustomValidators(v, database.DB)
	}

	routes.SetupRoutes(app)

	log.Fatal(app.Run(":3000"))
}
