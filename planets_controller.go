package main

import (
	"fmt"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	fmt.Printf("%s\n", "Retrieving planets")

	filter := bson.M{}

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
		fmt.Printf("%s\n", "Planets not found")
		c.Status(404).JSON(results)
	} else if err != nil {
		fmt.Printf("%s\n", "Can't retrieve planets due to internal error")
		c.Status(500).JSON(err)
	} else {
		fmt.Printf("%s\n", "Planets retrieved successfully")
		c.JSON(results)
	}
}

// Create render the planet create response as JSON
func (ctr *PlanetsControllers) Create(c *fiber.Ctx) {
	fmt.Printf("%s\n", "Creating planet")

	c.Accepts("application/json")

	resp, err := ctr.PlanetsClient.Create(c.Body())

	_, isWriteException := err.(mongo.WriteException)

	if err != nil && isWriteException && err.(mongo.WriteException).WriteErrors[0].Code == 11000 {
		fmt.Printf("%s\n", "Can't create planet due to conflict")
		c.Status(409).JSON(bson.M{"message": "'name' already exists"})
	} else if err != nil {
		fmt.Printf("%s\n", "Can't create planet due to internal error")
		c.Status(500).JSON(err)
	} else {
		fmt.Printf("%s\n", "PLanet created successfully")
		c.Status(201).JSON(resp)
	}
}

// Delete render the planet deletion response
func (ctr *PlanetsControllers) Delete(c *fiber.Ctx) {
	fmt.Printf("%s\n", "Deleting planet")

	deletedCount, err := ctr.PlanetsClient.Delete(c.Params("id"))

	if deletedCount == 0 {
		fmt.Printf("%s\n", "Planet not found")
		c.Status(404)
	} else if err != nil {
		fmt.Printf("%s\n", "Can't delete planet due to internal error")
		c.Status(500).JSON(err)
	} else {
		fmt.Printf("%s\n", "Planet deleted successfully")
		c.Status(204)
	}
}
