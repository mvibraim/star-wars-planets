package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const collection = "planets"

type PlanetsDatabase interface {
	Get(context.Context, interface{}) ([]Planet, error)
	Create(context.Context, *Planet) (interface{}, error)
	Delete(context.Context, string) (int64, error)
}

type planetsDatabase struct {
	db DatabaseHelper
}

func createPlanetsDatabase() PlanetsDatabase {
	clientHelper, _ := NewClient(config.MongoDBHost)
	databaseHelper := NewDatabase(config.MongoDBDatabase, clientHelper)
	planetsDatabase := NewPlanetsDatabase(databaseHelper)

	return planetsDatabase
}

func NewPlanetsDatabase(db DatabaseHelper) PlanetsDatabase {
	return &planetsDatabase{
		db: db,
	}
}

func (planetDB *planetsDatabase) Get(ctx context.Context, filter interface{}) ([]Planet, error) {
	planets := &[]Planet{}
	err := planetDB.db.Collection(collection).Find(ctx, filter).All(ctx, planets)

	if err != nil {
		return nil, err
	}

	return *planets, nil
}

func (planetDB *planetsDatabase) Create(ctx context.Context, planet *Planet) (interface{}, error) {
	res, err := planetDB.db.Collection(collection).InsertOne(ctx, planet)
	return res, err
}

func (planetDB *planetsDatabase) Delete(ctx context.Context, id string) (int64, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}

	res, err := planetDB.db.Collection(collection).DeleteOne(ctx, filter)

	return res, err
}
