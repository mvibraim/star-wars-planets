package main

import (
	"star-wars-planets/routes"

	"github.com/gofiber/compression"
	"github.com/gofiber/fiber"
	"github.com/gofiber/helmet"
	"github.com/gofiber/logger"
	"github.com/gofiber/recover"
)

const port = 8000

func main() {
	app := fiber.New()

	app.Use(compression.New())
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(helmet.New())

	routes.Planets(app)

	app.Listen(port)
}
