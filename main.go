package main

import (
	"log"
	"inventory/database"
	"inventory/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()

	database.Connect()

	routes.SetupRoutes(app)

	log.Fatal(app.Run(":3000"))
}
