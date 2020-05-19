package database

import (
	"context"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const dbName = "planets-db"
const collectionName = "planets"

type PlanetsClient struct{}

type Planet struct {
	Name    string `json:"name,omitempty"`
	Weather string `json:"weather,omitempty"`
	Terrain string `json:"terrain,omitempty"`
}

func (client PlanetsClient) Get(filter bson.M) ([]primitive.M, error) {
	collection, err := GetMongoDbCollection(dbName, collectionName)

	if err != nil {
		return nil, err
	}

	var results []primitive.M
	cur, err := collection.Find(context.Background(), filter)
	defer cur.Close(context.Background())

	if err != nil {
		return nil, err
	}

	cur.All(context.Background(), &results)

	if results == nil {
		return nil, nil
	}

	return results, err
}

func (client PlanetsClient) Create(body string) (map[string]string, error) {
	collection, err := GetMongoDbCollection(dbName, collectionName)

	if err != nil {
		return nil, err
	}

	var planet Planet
	json.Unmarshal([]byte(body), &planet)

	res, err := collection.InsertOne(context.Background(), planet)

	if err != nil {
		return nil, err
	}

	id, _ := res.InsertedID.(primitive.ObjectID)
	response := map[string]string{"id": id.Hex()}

	return response, err
}

func (client PlanetsClient) Delete(id string) (int64, error) {
	collection, err := GetMongoDbCollection(dbName, collectionName)

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
