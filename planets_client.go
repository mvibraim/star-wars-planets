package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const collectionName = "planets"

// PlanetsClient represents the client for Planet collection
type PlanetsClient struct{}

// Planet represents each record in Planets collection
type Planet struct {
	Name             string `json:"name,omitempty"`
	Climate          string `json:"climate,omitempty"`
	Terrain          string `json:"terrain,omitempty"`
	MovieAppearances int    `bson:"movie_appearances" json:"-"`
}

// Get return planets from database
func (client PlanetsClient) Get(filter bson.M) ([]Planet, error) {
	collection, err := GetMongoDbCollection(config.MongoDBDatabase, collectionName)

	if err != nil {
		return nil, err
	}

	var results []Planet
	cur, err := collection.Find(context.Background(), filter)
	defer cur.Close(context.Background())

	if err != nil {
		return nil, err
	}

	cur.All(context.Background(), &results)

	if results == nil {
		return []Planet{}, nil
	}

	return results, err
}

// Create insert a planet in database
func (client PlanetsClient) Create(body string) (map[string]string, error) {
	collection, err := GetMongoDbCollection(config.MongoDBDatabase, collectionName)

	if err != nil {
		return nil, err
	}

	var planet Planet
	json.Unmarshal([]byte(body), &planet)

	conn := getRedisConn()
	filmsCount, _ := getCache(conn, strings.ToLower(planet.Name))
	conn.Close()

	if filmsCount == -1 {
		planet.MovieAppearances = 0
	} else {
		planet.MovieAppearances = filmsCount
	}

	fmt.Printf("%s has %d movie appearances\n", planet.Name, planet.MovieAppearances)

	res, err := collection.InsertOne(context.Background(), planet)

	if err != nil {
		return nil, err
	}

	id, _ := res.InsertedID.(primitive.ObjectID)
	response := map[string]string{"id": id.Hex()}

	return response, err
}

// Delete deletes a planet in database
func (client PlanetsClient) Delete(id string) (int64, error) {
	collection, err := GetMongoDbCollection(config.MongoDBDatabase, collectionName)

	if err != nil {
		return -1, err
	}

	objID, _ := primitive.ObjectIDFromHex(id)
	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})

	if err != nil {
		return -1, err
	}

	return result.DeletedCount, nil
}
