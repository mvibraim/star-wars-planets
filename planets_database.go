package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const collection = "planets"

type PlanetsDatabase interface {
	Get(context.Context, interface{}) ([]Planet, error)
	Create(context.Context, *Planet) (*mongo.InsertOneResult, error)
	Delete(context.Context, string) (int64, error)
	CreateIndexes(context.Context, mongo.IndexModel) (string, error)
}

type planetsDatabase struct {
	db DatabaseHelper
}

func setupPlanetsDatabase() PlanetsDatabase {
	clientHelper, _ := NewClient(config.MongoDBHost)
	databaseHelper := NewDatabase(config.MongoDBDatabase, clientHelper)
	planetsDatabase := NewPlanetsDatabase(databaseHelper)

	indexModel := mongo.IndexModel{
		Keys:    bsonx.Doc{{"name", bsonx.Int32(1)}},
		Options: options.Index().SetUnique(true),
	}

	planetsDatabase.CreateIndexes(context.Background(), indexModel)

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

func (planetDB *planetsDatabase) Create(ctx context.Context, planet *Planet) (*mongo.InsertOneResult, error) {
	res, err := planetDB.db.Collection(collection).InsertOne(ctx, planet)
	return res, err
}

func (planetDB *planetsDatabase) Delete(ctx context.Context, id string) (int64, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}

	res, err := planetDB.db.Collection(collection).DeleteOne(ctx, filter)

	return res, err
}

func (planetDB *planetsDatabase) CreateIndexes(ctx context.Context, indexModel mongo.IndexModel) (string, error) {
	res, err := planetDB.db.Collection(collection).Indexes().CreateOne(ctx, indexModel)
	return res, err
}
