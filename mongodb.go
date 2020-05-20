package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoConnection, errConnection = GetMongoDbConnection()

// GetMongoDbConnection gets the MongoDB connection
func GetMongoDbConnection() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(config.MongoDBHost)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())

	if err != nil {
		log.Fatal(err)
	}

	return client, nil
}

// GetMongoDbCollection gets the MongoDB collection
func GetMongoDbCollection(DbName string, CollectionName string) (*mongo.Collection, error) {
	if errConnection != nil {
		return nil, errConnection
	}

	collection := mongoConnection.Database(DbName).Collection(CollectionName)

	return collection, nil
}
