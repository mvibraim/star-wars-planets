package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	planets, err := planetsDB.Get(context.Background(), filter)
	return planets, err
}

// Create insert a planet in database
func (client PlanetsClient) Create(body string) (map[string]string, error) {
	var planet Planet
	json.Unmarshal([]byte(body), &planet)

	conn := getRedisConn()
	movieAppearances, _ := getCache(conn, strings.ToLower(planet.Name))
	conn.Close()

	if movieAppearances == -1 {
		planet.MovieAppearances = 0
	} else {
		planet.MovieAppearances = movieAppearances
	}

	fmt.Printf("%s has %d movie appearances\n", planet.Name, planet.MovieAppearances)

	insertedID, err := planetsDB.Create(context.Background(), &planet)

	if err != nil {
		return nil, err
	}

	id, _ := insertedID.(primitive.ObjectID)
	response := map[string]string{"id": id.Hex()}

	return response, err
}

// Delete deletes a planet in database
func (client PlanetsClient) Delete(id string) (int64, error) {
	deletedCount, err := planetsDB.Delete(context.Background(), id)

	if err != nil {
		return -1, err
	}

	return deletedCount, nil
}
