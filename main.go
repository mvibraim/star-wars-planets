package main

import (
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

	app.Post("/planets", create)
	app.Get("/planets", index)
	app.Delete("/planets/:id", delete)

	app.Listen(port)
}
