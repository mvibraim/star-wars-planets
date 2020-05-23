package main

import (
	"github.com/gofiber/compression"
	"github.com/gofiber/fiber"
	"github.com/gofiber/helmet"
	"github.com/gofiber/logger"
	"github.com/gofiber/recover"
	"go.mongodb.org/mongo-driver/mongo"
)

var planetsDB PlanetsDatabase
var config Config = parseConfig()
var mongoBDClient *mongo.Client
var mongoConnectionError error
var isTesting bool = false

func main() {
	if !isTesting {
		planetsDB = setupPlanetsDatabase()
		cacheMovieAppearancesByName()
	}

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
