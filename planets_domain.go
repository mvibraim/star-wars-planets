package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlanetsDBHelper interface {
	Get(context.Context, interface{}) ([]Planet, error)
	Create(context.Context, *Planet) (*mongo.InsertOneResult, error)
	Delete(context.Context, string) (int64, error)
}

// PlanetsDomain represents the client for Planet collection
type PlanetsDomain struct {
	PlanetsDB    PlanetsDBHelper
	PlanetsCache PlanetsCacheHelper
}

// Planet represents each record in Planets collection
type Planet struct {
	Name             string `json:"name,omitempty" validate:"required"`
	Climate          string `json:"climate,omitempty" validate:"required"`
	Terrain          string `json:"terrain,omitempty" validate:"required"`
	MovieAppearances int    `bson:"movie_appearances" json:"-"`
}

func CreatePlanetsDomain() *PlanetsDomain {
	return &PlanetsDomain{
		PlanetsDB:    CreatePlanetsDB(),
		PlanetsCache: CreatePlanetsCache(),
	}
}

// Get return planets from database
func (domain *PlanetsDomain) Get(filter bson.M) ([]Planet, error) {
	planets, err := domain.PlanetsDB.Get(context.Background(), filter)
	return planets, err
}

// Create insert a planet in database
func (domain *PlanetsDomain) Create(body string) (map[string]string, error) {
	var planet Planet
	json.Unmarshal([]byte(body), &planet)

	v := validator.New()
	validationErrors := v.Struct(planet)

	if validationErrors != nil {
		return nil, validationErrors
	}

	movieAppearances, _ := domain.PlanetsCache.getCache(planet.Name)

	if movieAppearances == -1 {
		planet.MovieAppearances = 0
	} else {
		planet.MovieAppearances = movieAppearances
	}

	fmt.Printf("%s has %d movie appearances\n", planet.Name, planet.MovieAppearances)

	res, err := domain.PlanetsDB.Create(context.Background(), &planet)

	if err != nil {
		return nil, err
	}

	id, _ := res.InsertedID.(primitive.ObjectID)
	response := map[string]string{"id": id.Hex()}

	return response, err
}

// Delete deletes a planet in database
func (domain *PlanetsDomain) Delete(id string) (int64, error) {
	deletedCount, err := domain.PlanetsDB.Delete(context.Background(), id)

	if err != nil {
		return -1, err
	}

	return deletedCount, nil
}
