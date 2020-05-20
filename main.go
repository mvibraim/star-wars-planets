package main

import (
	"github.com/gofiber/compression"
	"github.com/gofiber/fiber"
	"github.com/gofiber/helmet"
	"github.com/gofiber/logger"
	"github.com/gofiber/recover"
)

const port = 3000

func main() {
	cacheFilmsCountByName()

	app := fiber.New()

	app.Use(compression.New())
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(helmet.New())

	PlanetsRoutes(app)

	app.Listen(port)
}

// PlanetsRoutes setup all the routes related to planets
func PlanetsRoutes(app *fiber.App) {
	planetsClient := new(PlanetsClient)
	ctr := PlanetsControllers{}
	ctr.PlanetsClient = planetsClient

	planets := app.Group("v1/planets")

	planets.Post("/", ctr.Create)
	planets.Get("/", ctr.Index)
	planets.Delete("/:id", ctr.Delete)
}
