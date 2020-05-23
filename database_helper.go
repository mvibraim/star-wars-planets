package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseHelper interface {
	Collection(name string) CollectionHelper
	Client() ClientHelper
}

type CollectionHelper interface {
	Find(context.Context, interface{}) CursorHelper
	InsertOne(context.Context, interface{}) (*mongo.InsertOneResult, error)
	DeleteOne(ctx context.Context, filter interface{}) (int64, error)
	Indexes() IndexViewHelper
}

type ClientHelper interface {
	Database(string) DatabaseHelper
	Connect() error
	StartSession() (mongo.Session, error)
}

type CursorHelper interface {
	All(ctx context.Context, results interface{}) error
}

type IndexViewHelper interface {
	CreateOne(ctx context.Context, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error)
}

type MongoClient struct {
	cl *mongo.Client
}

type MongoDatabase struct {
	db *mongo.Database
}

type MongoCollection struct {
	coll *mongo.Collection
}

type MongoCursor struct {
	c *mongo.Cursor
}

type MongoIndexView struct {
	iv mongo.IndexView
}

type mongoSession struct {
	mongo.Session
}

func CreateClient(databaseHost string) (ClientHelper, error) {
	c, err := mongo.Connect(context.Background(), options.Client().ApplyURI(databaseHost))
	return &MongoClient{cl: c}, err
}

func CreateDatabase(databaseName string, client ClientHelper) DatabaseHelper {
	return client.Database(databaseName)
}

func (mc *MongoClient) Database(dbName string) DatabaseHelper {
	db := mc.cl.Database(dbName)
	return &MongoDatabase{db: db}
}

func (mc *MongoClient) StartSession() (mongo.Session, error) {
	session, err := mc.cl.StartSession()
	return &mongoSession{session}, err
}

func (mc *MongoClient) Connect() error {
	return mc.cl.Connect(nil)
}

func (md *MongoDatabase) Collection(colName string) CollectionHelper {
	collection := md.db.Collection(colName)
	return &MongoCollection{coll: collection}
}

func (md *MongoDatabase) Client() ClientHelper {
	client := md.db.Client()
	return &MongoClient{cl: client}
}

func (mc *MongoCollection) Find(ctx context.Context, filter interface{}) CursorHelper {
	cursor, _ := mc.coll.Find(ctx, filter)
	return &MongoCursor{c: cursor}
}

func (mc *MongoCollection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	res, err := mc.coll.InsertOne(ctx, document)
	return res, err
}

func (mc *MongoCollection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	count, err := mc.coll.DeleteOne(ctx, filter)
	return count.DeletedCount, err
}

func (mc *MongoCollection) Indexes() IndexViewHelper {
	indexView := mc.coll.Indexes()
	return &MongoIndexView{iv: indexView}
}

func (c *MongoCursor) All(ctx context.Context, results interface{}) error {
	return c.c.All(ctx, results)
}

func (iv *MongoIndexView) CreateOne(ctx context.Context, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error) {
	return iv.iv.CreateOne(ctx, model, opts...)
}
