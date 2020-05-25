package main

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlanetsDomainMock struct {
	mock.Mock
}

func (domain *PlanetsDomainMock) Get(filter bson.M) ([]Planet, error) {
	args := domain.Called(filter)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Planet), args.Error(1)
}

func (domain *PlanetsDomainMock) Create(body string) (map[string]string, error) {
	args := domain.Called(body)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(map[string]string), args.Error(1)
}

func (domain *PlanetsDomainMock) Delete(id string) (int64, error) {
	args := domain.Called(id)
	return args.Get(0).(int64), args.Error(1)
}

type PlanetsCacheMock struct {
	mock.Mock
}

func (cache *PlanetsCacheMock) SetCache(name string, movieAppearances int) error {
	args := cache.Called(name, movieAppearances)

	if args.Get(0) == nil {
		return nil
	}

	return args.Error(0)
}

func (cache *PlanetsCacheMock) GetCache(name string) (int, error) {
	args := cache.Called(name)

	if args.Get(1) == nil {
		return args.Int(0), nil
	}

	return args.Int(0), args.Error(1)
}

type PlanetsHttpClientMock struct {
	mock.Mock
}

func (client *PlanetsHttpClientMock) Get(url string) (*http.Response, error) {
	args := client.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}

type PlanetsDBMock struct {
	mock.Mock
}

func (db *PlanetsDBMock) Get(ctx context.Context, filter interface{}) ([]Planet, error) {
	args := db.Called(ctx, filter)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).([]Planet), args.Error(1)
	}
}

func (db *PlanetsDBMock) Create(ctx context.Context, planet *Planet) (*mongo.InsertOneResult, error) {
	args := db.Called(ctx, planet)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (db *PlanetsDBMock) Delete(ctx context.Context, id string) (int64, error) {
	args := db.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}
