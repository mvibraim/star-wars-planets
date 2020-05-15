package main

import (
	"github.com/gofiber/compression"
	"github.com/gofiber/fiber"
)

const port = 8000

func main() {
	app := fiber.New()

	app.Use(compression.New())

	app.Post("/planets", create)

	app.Get("/planets", index)

	app.Delete("/planets/:id", delete)

	app.Listen(port)
}
