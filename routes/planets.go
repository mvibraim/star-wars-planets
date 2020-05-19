package routes

import (
	"star-wars-planets/controllers"
	"star-wars-planets/database"

	"github.com/gofiber/fiber"
)

func Planets(app *fiber.App) {
	planetsClient := new(database.PlanetsClient)
	ctr := controllers.Controllers{}
	ctr.PlanetsClient = planetsClient

	planets := app.Group("v1/planets")

	planets.Post("/", ctr.Create)
	planets.Get("/", ctr.Index)
	planets.Delete("/:id", ctr.Delete)
}
