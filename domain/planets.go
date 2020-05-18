package domain

import (
	"context"
	"encoding/json"
	"star-wars-planets/database"
	"star-wars-planets/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const dbName = "planets-db"
const collectionName = "planets"

type Planets struct{}

func (planets Planets) Get(filter bson.M) ([]bson.M, error) {
	collection, err := database.GetMongoDbCollection(dbName, collectionName)

	if err != nil {
		return nil, err
	}

	var results []bson.M
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

func (planets Planets) Create(body string) (map[string]string, error) {
	collection, err := database.GetMongoDbCollection(dbName, collectionName)

	if err != nil {
		return nil, err
	}

	var planet models.Planet
	json.Unmarshal([]byte(body), &planet)

	res, err := collection.InsertOne(context.Background(), planet)

	if err != nil {
		return nil, err
	}

	id, _ := res.InsertedID.(primitive.ObjectID)
	response := map[string]string{"id": id.Hex()}

	return response, err
}

func (planets Planets) Delete(id string) (int64, error) {
	collection, err := database.GetMongoDbCollection(dbName, collectionName)

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
