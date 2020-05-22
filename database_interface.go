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

type CursorHelper interface {
	All(ctx context.Context, results interface{}) error
}

type ClientHelper interface {
	Database(string) DatabaseHelper
	Connect() error
	StartSession() (mongo.Session, error)
}

type IndexViewHelper interface {
	CreateOne(ctx context.Context, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error)
}

type mongoClient struct {
	cl *mongo.Client
}

type mongoDatabase struct {
	db *mongo.Database
}

type mongoCollection struct {
	coll *mongo.Collection
}

type mongoCursor struct {
	c *mongo.Cursor
}

type mongoIndexView struct {
	iv mongo.IndexView
}

type mongoSession struct {
	mongo.Session
}

func NewClient(databaseHost string) (ClientHelper, error) {
	c, err := mongo.Connect(context.Background(), options.Client().ApplyURI(databaseHost))
	return &mongoClient{cl: c}, err
}

func NewDatabase(databaseName string, client ClientHelper) DatabaseHelper {
	return client.Database(databaseName)
}

func (mc *mongoClient) Database(dbName string) DatabaseHelper {
	db := mc.cl.Database(dbName)
	return &mongoDatabase{db: db}
}

func (mc *mongoClient) StartSession() (mongo.Session, error) {
	session, err := mc.cl.StartSession()
	return &mongoSession{session}, err
}

func (mc *mongoClient) Connect() error {
	return mc.cl.Connect(nil)
}

func (md *mongoDatabase) Collection(colName string) CollectionHelper {
	collection := md.db.Collection(colName)
	return &mongoCollection{coll: collection}
}

func (md *mongoDatabase) Client() ClientHelper {
	client := md.db.Client()
	return &mongoClient{cl: client}
}

func (mc *mongoCollection) Find(ctx context.Context, filter interface{}) CursorHelper {
	cursor, _ := mc.coll.Find(ctx, filter)
	return &mongoCursor{c: cursor}
}

func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	res, err := mc.coll.InsertOne(ctx, document)
	return res, err
}

func (mc *mongoCollection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	count, err := mc.coll.DeleteOne(ctx, filter)
	return count.DeletedCount, err
}

func (mc *mongoCollection) Indexes() IndexViewHelper {
	indexView := mc.coll.Indexes()
	return &mongoIndexView{iv: indexView}
}

func (c *mongoCursor) All(ctx context.Context, results interface{}) error {
	return c.c.All(ctx, results)
}

func (iv *mongoIndexView) CreateOne(ctx context.Context, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error) {
	return iv.iv.CreateOne(ctx, model, opts...)
}
