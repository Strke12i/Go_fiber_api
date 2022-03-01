package main

import (
	"freecode/database"
	"freecode/routes"
	"github.com/gofiber/fiber"
)

func main() {

	database.Connect()

	app := fiber.New()
	
	routes.Setup(app)

	app.Listen(":8000")
}
