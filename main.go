package main

import (
	"context"
	"encoding/json"

	"github.com/gofiber/compression"
	"github.com/gofiber/fiber"
	"github.com/gofiber/helmet"
	"github.com/gofiber/logger"
	"github.com/gofiber/recover"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const dbName = "planets-db"
const collectionName = "planets"
const port = 8000

func index(c *fiber.Ctx) {
	collection, err := getMongoDbCollection(dbName, collectionName)

	if err != nil {
		c.Status(500).JSON(err)
		return
	}

	var filter bson.M = bson.M{}

	if c.Query("id") != "" {
		id := c.Query("id")
		objID, _ := primitive.ObjectIDFromHex(id)
		filter = bson.M{"_id": objID}
	} else if c.Query("name") != "" {
		name := c.Query("name")
		filter = bson.M{"name": name}
	}

	var results []bson.M
	cur, err := collection.Find(context.Background(), filter)
	defer cur.Close(context.Background())

	if err != nil {
		c.Status(500).JSON(err)
		return
	}

	cur.All(context.Background(), &results)

	if results == nil {
		c.SendStatus(404)
		return
	}

	c.JSON(results)
}

func create(c *fiber.Ctx) {
	c.Accepts("application/json")

	collection, err := getMongoDbCollection(dbName, collectionName)

	if err != nil {
		c.Status(500).JSON(err)
		return
	}

	var planet Planet
	json.Unmarshal([]byte(c.Body()), &planet)

	res, err := collection.InsertOne(context.Background(), planet)

	if err != nil {
		c.Status(500).JSON(err)
		return
	}

	id, _ := res.InsertedID.(primitive.ObjectID)
	response := map[string]string{"id": id.Hex()}

	c.Status(201).JSON(response)
}

func delete(c *fiber.Ctx) {
	collection, err := getMongoDbCollection(dbName, collectionName)

	if err != nil {
		c.Status(500).JSON(err)
		return
	}

	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objID})

	if err != nil {
		c.Status(500).JSON(err)
		return
	}

	c.Status(204)
}

func main() {
	app := fiber.New()

	app.Use(compression.New())
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(helmet.New())

	app.Post("v1/planets", create)
	app.Get("v1/planets", index)
	app.Delete("v1/planets/:id", delete)

	app.Listen(port)
}
