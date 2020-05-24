package main

import (
	"github.com/gofiber/compression"
	"github.com/gofiber/fiber"
	"github.com/gofiber/helmet"
	"github.com/gofiber/logger"
	"github.com/gofiber/recover"
)

var config Config = parseConfig()

func main() {
	CreateSwapiClient().CacheMovieAppearancesByName()

	app := fiber.New()

	app.Use(compression.New())
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(helmet.New())

	PlanetsRoutes(app)

	app.Listen(config.Port)
}

// PlanetsRoutes setup all the routes related to planets
func PlanetsRoutes(app *fiber.App) {
	ctr := CreatePlanetsController()

	planets := app.Group("v1/planets")

	planets.Post("/", ctr.Create)
	planets.Get("/", ctr.Index)
	planets.Delete("/:id", ctr.Delete)
}
