package main

import (
	"fmt"

	"github.com/gofiber/compression"
	"github.com/gofiber/fiber"
	"github.com/gofiber/helmet"
	"github.com/gofiber/logger"
	"github.com/gofiber/recover"
	"go.mongodb.org/mongo-driver/mongo"
)

var config Config
var mongoClient *mongo.Client
var mongoConnectionError error
var isTesting bool = false

func main() {
	config = parseConfig()

	fmt.Printf("%t\n", isTesting)

	if !isTesting {
		mongoClient, mongoConnectionError = GetMongoDbConnection()
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
	planetsClient := new(PlanetsClient)
	ctr := PlanetsControllers{}
	ctr.PlanetsClient = planetsClient

	planets := app.Group("v1/planets")

	planets.Post("/", ctr.Create)
	planets.Get("/", ctr.Index)
	planets.Delete("/:id", ctr.Delete)
}
