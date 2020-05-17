package routes

import (
	"star-wars-planets/controllers"

	"github.com/gofiber/fiber"
)

func Planets(app *fiber.App) {
	planets := app.Group("v1/planets")
	planets.Post("/", controllers.Create)
	planets.Get("/", controllers.Index)
	planets.Delete("/:id", controllers.Delete)
}
