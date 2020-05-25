package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestGetPlanetsFromDatabaseSuccessfully(t *testing.T) {
	filter := bson.M{}
	response := []Planet{}
	planetsDBMock := new(PlanetsDBMock)
	planetsDBMock.On("Get", mock.Anything, filter).Return(response, nil)

	planetsDomain := PlanetsDomain{
		PlanetsDB: planetsDBMock,
	}

	resp, err := planetsDomain.Get(filter)

	planetsDBMock.AssertExpectations(t)

	assert := assert.New(t)

	assert.Equal(nil, err, "they should be equal")
	assert.Equal(response, resp, "they should be equal")
}

func TestDontGetPlanetsFromDatabaseDueToError(t *testing.T) {
	filter := bson.M{}
	planetsDBMock := new(PlanetsDBMock)
	planetsDBMock.On("Get", mock.Anything, filter).Return(nil, errors.New(""))

	planetsDomain := PlanetsDomain{
		PlanetsDB: planetsDBMock,
	}

	_, err := planetsDomain.Get(filter)

	planetsDBMock.AssertExpectations(t)

	assert := assert.New(t)

	assert.NotEqual(nil, err, "they shouldn't be equal")
}

func TestDeletePlanetsFromDatabaseSuccessfully(t *testing.T) {
	id := "id"
	response := int64(1)
	planetsDBMock := new(PlanetsDBMock)
	planetsDBMock.On("Delete", mock.Anything, id).Return(response, nil)

	planetsDomain := PlanetsDomain{
		PlanetsDB: planetsDBMock,
	}

	resp, err := planetsDomain.Delete(id)

	planetsDBMock.AssertExpectations(t)

	assert := assert.New(t)

	assert.Equal(nil, err, "they should be equal")
	assert.Equal(response, resp, "they should be equal")
}

func TestDontDeletePlanetsFromDatabaseDueToError(t *testing.T) {
	id := "id"
	response := int64(-1)
	planetsDBMock := new(PlanetsDBMock)
	planetsDBMock.On("Delete", mock.Anything, id).Return(response, errors.New(""))

	planetsDomain := PlanetsDomain{
		PlanetsDB: planetsDBMock,
	}

	resp, err := planetsDomain.Delete(id)

	planetsDBMock.AssertExpectations(t)

	assert := assert.New(t)

	assert.NotEqual(nil, err, "they shouldn't be equal")
	assert.Equal(response, resp, "they should be equal")
}

func TestCreatePlanetsToDatabaseSuccessfully(t *testing.T) {
	body := `{"name":"Tatooine","climate":"hot","terrain":"dust"}`
	planet := Planet{Name: "Tatooine", Climate: "hot", Terrain: "dust", MovieAppearances: 5}
	objectID := primitive.NewObjectID()
	dbCreateResponse := mongo.InsertOneResult{InsertedID: objectID}
	response := map[string]string{"id": objectID.Hex()}

	planetsDBMock := new(PlanetsDBMock)
	planetsDBMock.On("Create", mock.Anything, &planet).Return(&dbCreateResponse, nil)

	planetsCacheMock := new(PlanetsCacheMock)
	planetsCacheMock.On("GetCache", "Tatooine").Return(5, nil)

	planetsDomain := PlanetsDomain{
		PlanetsCache: planetsCacheMock,
		PlanetsDB:    planetsDBMock,
	}

	resp, err := planetsDomain.Create(body)

	planetsCacheMock.AssertExpectations(t)
	planetsDBMock.AssertExpectations(t)

	assert := assert.New(t)

	assert.Equal(nil, err, "they should be equal")
	assert.Equal(response, resp, "they should be equal")
}

func TestDontCreatePlanetsToDatabaseDueToValidationErrors(t *testing.T) {
	body := `{}`
	planetsDomain := PlanetsDomain{}

	_, err := planetsDomain.Create(body)

	assert := assert.New(t)

	assert.NotEqual(nil, err, "they shouldn't be equal")
}

func TestDontCreatePlanetsToDatabaseDueToCacheError(t *testing.T) {
	body := `{"name":"Tatooine","climate":"hot","terrain":"dust"}`

	planetsCacheMock := new(PlanetsCacheMock)
	planetsCacheMock.On("GetCache", "Tatooine").Return(-1, errors.New(""))

	planetsDomain := PlanetsDomain{
		PlanetsCache: planetsCacheMock,
	}

	resp, err := planetsDomain.Create(body)

	planetsCacheMock.AssertExpectations(t)

	assert := assert.New(t)

	assert.NotEqual(nil, err, "they shouldn't be equal")
	assert.Equal(map[string]string(nil), resp, "they should be equal")
}

func TestDontCreatePlanetsToDatabaseDueToCreateError(t *testing.T) {
	body := `{"name":"Tatooine","climate":"hot","terrain":"dust"}`
	planet := Planet{Name: "Tatooine", Climate: "hot", Terrain: "dust", MovieAppearances: 5}
	dbCreateResponse := mongo.InsertOneResult{}

	planetsCacheMock := new(PlanetsCacheMock)
	planetsCacheMock.On("GetCache", "Tatooine").Return(5, nil)

	planetsDBMock := new(PlanetsDBMock)
	planetsDBMock.On("Create", mock.Anything, &planet).Return(&dbCreateResponse, errors.New(""))

	planetsDomain := PlanetsDomain{
		PlanetsCache: planetsCacheMock,
		PlanetsDB:    planetsDBMock,
	}

	_, err := planetsDomain.Create(body)

	planetsCacheMock.AssertExpectations(t)
	planetsDBMock.AssertExpectations(t)

	assert := assert.New(t)

	assert.NotEqual(nil, err, "they shouldn't be equal")
}
