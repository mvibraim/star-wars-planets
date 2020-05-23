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

type PlanetsDatabase struct {
	db DatabaseHelper
}

func CreateDatabaseHelper() DatabaseHelper {
	clientHelper, _ := CreateClient(config.MongoDBHost)
	databaseHelper := CreateDatabase(config.MongoDBDatabase, clientHelper)

	return databaseHelper
}

func CreatePlanetsDB() *PlanetsDatabase {
	indexModel := mongo.IndexModel{
		Keys:    bsonx.Doc{{"name", bsonx.Int32(1)}},
		Options: options.Index().SetUnique(true),
	}

	planetsDatabase := &PlanetsDatabase{
		db: CreateDatabaseHelper(),
	}

	planetsDatabase.CreateIndexes(context.Background(), indexModel)

	return planetsDatabase
}

func (planetsDB *PlanetsDatabase) Get(ctx context.Context, filter interface{}) ([]Planet, error) {
	planets := &[]Planet{}
	err := planetsDB.db.Collection(collection).Find(ctx, filter).All(ctx, planets)

	if err != nil {
		return nil, err
	}

	return *planets, nil
}

func (planetsDB *PlanetsDatabase) Create(ctx context.Context, planet *Planet) (*mongo.InsertOneResult, error) {
	res, err := planetsDB.db.Collection(collection).InsertOne(ctx, planet)
	return res, err
}

func (planetsDB *PlanetsDatabase) Delete(ctx context.Context, id string) (int64, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}

	res, err := planetsDB.db.Collection(collection).DeleteOne(ctx, filter)

	return res, err
}

func (planetsDB *PlanetsDatabase) CreateIndexes(ctx context.Context, indexModel mongo.IndexModel) (string, error) {
	res, err := planetsDB.db.Collection(collection).Indexes().CreateOne(ctx, indexModel)
	return res, err
}
