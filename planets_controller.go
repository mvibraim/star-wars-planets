package main

import (
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PlanetsControllers represents the planet controller structure
type PlanetsControllers struct {
	PlanetsClient interface {
		Get(filter bson.M) ([]Planet, error)
		Create(body string) (map[string]string, error)
		Delete(id string) (int64, error)
	}
}

// Index render the returned planets as JSON
func (ctr *PlanetsControllers) Index(c *fiber.Ctx) {
	var filter bson.M = bson.M{}

	if c.Query("id") != "" {
		id := c.Query("id")
		objID, _ := primitive.ObjectIDFromHex(id)
		filter = bson.M{"_id": objID}
	} else if c.Query("name") != "" {
		name := c.Query("name")
		filter = bson.M{"name": name}
	}

	results, err := ctr.PlanetsClient.Get(filter)

	if len(results) == 0 && err == nil {
		c.Status(404).JSON(results)
	} else if err != nil {
		c.Status(500).JSON(err)
	} else {
		c.JSON(results)
	}
}

// Create render the planet create response as JSON
func (ctr *PlanetsControllers) Create(c *fiber.Ctx) {
	c.Accepts("application/json")

	resp, err := ctr.PlanetsClient.Create(c.Body())

	if err != nil {
		c.Status(500).JSON(err)
	} else {
		c.Status(201).JSON(resp)
	}
}

// Delete render the planet deletion response
func (ctr *PlanetsControllers) Delete(c *fiber.Ctx) {
	deletedCount, err := ctr.PlanetsClient.Delete(c.Params("id"))

	if deletedCount == 0 {
		c.Status(404)
	} else if err != nil {
		c.Status(500).JSON(err)
	} else {
		c.Status(204)
	}
}
