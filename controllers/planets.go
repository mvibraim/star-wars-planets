package controllers

import (
	"star-wars-planets/domain"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Index(c *fiber.Ctx) {
	var filter bson.M = bson.M{}

	if c.Query("id") != "" {
		id := c.Query("id")
		objID, _ := primitive.ObjectIDFromHex(id)
		filter = bson.M{"_id": objID}
	} else if c.Query("name") != "" {
		name := c.Query("name")
		filter = bson.M{"name": name}
	}

	results, err := domain.GetPlanets(filter)

	if results == nil && err == nil {
		c.SendStatus(404)
	} else if err != nil {
		c.Status(500).JSON(err)
	} else {
		c.JSON(results)
	}
}

func Create(c *fiber.Ctx) {
	c.Accepts("application/json")

	response, err := domain.CreatePlanet(c.Body())

	if err != nil {
		c.Status(500).JSON(err)
	} else {
		c.Status(201).JSON(response)
	}
}

func Delete(c *fiber.Ctx) {
	deletedCount, err := domain.DeletePlanet(c.Params("id"))

	if deletedCount == 0 {
		c.Status(404)
	} else if err != nil {
		c.Status(500).JSON(err)
	} else {
		c.Status(204)
	}
}
