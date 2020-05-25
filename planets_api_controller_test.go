package main

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestIntegrationGetPlanetsSuccessfully(t *testing.T) {
	filter := bson.M{}
	response := []Planet{{Name: "bar", Terrain: "world", Climate: "cold"}}

	planetsDBMock := new(PlanetsDBMock)
	planetsDBMock.On("Get", mock.Anything, filter).Return(response, nil)

	ctr := PlanetsController{
		PlanetsDomain: &PlanetsDomain{
			PlanetsDB: planetsDBMock,
		},
	}

	app := fiber.New()
	app.Get("/v1/planets", ctr.Index)

	req := httptest.NewRequest("GET", "/v1/planets", nil)
	resp, _ := app.Test(req)

	planetsDBMock.AssertExpectations(t)

	assert := assert.New(t)
	assert.Equal(200, resp.StatusCode, "they should be equal")
}

func TestIntegrationGetPlanetsWithNameFilterSuccessfully(t *testing.T) {
	filter := bson.M{"name": "name"}
	response := []Planet{{Name: "bar", Terrain: "world", Climate: "cold"}}

	planetsDBMock := new(PlanetsDBMock)
	planetsDBMock.On("Get", mock.Anything, filter).Return(response, nil)

	ctr := PlanetsController{
		PlanetsDomain: &PlanetsDomain{
			PlanetsDB: planetsDBMock,
		},
	}

	app := fiber.New()
	app.Get("/v1/planets", ctr.Index)

	req := httptest.NewRequest("GET", "/v1/planets?name=name", nil)
	resp, _ := app.Test(req)

	planetsDBMock.AssertExpectations(t)

	assert := assert.New(t)
	assert.Equal(200, resp.StatusCode, "they should be equal")
}

func TestIntegrationGetPlanetsWithIdFilterSuccessfully(t *testing.T) {
	id := primitive.NewObjectID()
	filter := bson.M{"_id": id}
	response := []Planet{{Name: "bar", Terrain: "world", Climate: "cold"}}

	planetsDBMock := new(PlanetsDBMock)
	planetsDBMock.On("Get", mock.Anything, filter).Return(response, nil)

	ctr := PlanetsController{
		PlanetsDomain: &PlanetsDomain{
			PlanetsDB: planetsDBMock,
		},
	}
	app := fiber.New()
	app.Get("/v1/planets", ctr.Index)

	req := httptest.NewRequest("GET", "/v1/planets?id="+id.Hex(), nil)
	resp, _ := app.Test(req)

	planetsDBMock.AssertExpectations(t)

	assert := assert.New(t)
	assert.Equal(200, resp.StatusCode, "they should be equal")
}

func TestIntegrationDontGetPlanetsDueToNotFound(t *testing.T) {
	filter := bson.M{}
	response := []Planet{}

	planetsDBMock := new(PlanetsDBMock)
	planetsDBMock.On("Get", mock.Anything, filter).Return(response, nil)

	ctr := PlanetsController{
		PlanetsDomain: &PlanetsDomain{
			PlanetsDB: planetsDBMock,
		},
	}

	app := fiber.New()
	app.Get("/v1/planets", ctr.Index)

	req := httptest.NewRequest("GET", "/v1/planets", nil)
	resp, _ := app.Test(req)

	planetsDBMock.AssertExpectations(t)

	assert := assert.New(t)
	assert.Equal(404, resp.StatusCode, "they should be equal")
}

func TestIntegrationDontGetPlanetsDueToInternalError(t *testing.T) {
	filter := bson.M{}
	response := []Planet(nil)

	planetsDBMock := new(PlanetsDBMock)
	planetsDBMock.On("Get", mock.Anything, filter).Return(response, errors.New(""))

	ctr := PlanetsController{
		PlanetsDomain: &PlanetsDomain{
			PlanetsDB: planetsDBMock,
		},
	}

	app := fiber.New()
	app.Get("/v1/planets", ctr.Index)

	req := httptest.NewRequest("GET", "/v1/planets", nil)
	resp, _ := app.Test(req)

	planetsDBMock.AssertExpectations(t)

	assert := assert.New(t)
	assert.Equal(500, resp.StatusCode, "they should be equal")
}

func TestIntegrationDeletePlanetSuccsessfully(t *testing.T) {
	id := "id"
	response := int64(1)

	planetsDBMock := new(PlanetsDBMock)
	planetsDBMock.On("Delete", mock.Anything, id).Return(response, nil)

	ctr := PlanetsController{
		PlanetsDomain: &PlanetsDomain{
			PlanetsDB: planetsDBMock,
		},
	}

	app := fiber.New()
	app.Delete("/v1/planets/:id", ctr.Delete)

	req := httptest.NewRequest("DELETE", "/v1/planets/"+id, nil)
	resp, _ := app.Test(req)

	planetsDBMock.AssertExpectations(t)

	assert := assert.New(t)
	assert.Equal(204, resp.StatusCode, "they should be equal")
}

func TestIntegrationDontDeletePlanetDueToNotFound(t *testing.T) {
	id := "id"
	response := int64(0)

	planetsDBMock := new(PlanetsDBMock)
	planetsDBMock.On("Delete", mock.Anything, id).Return(response, nil)

	ctr := PlanetsController{
		PlanetsDomain: &PlanetsDomain{
			PlanetsDB: planetsDBMock,
		},
	}

	app := fiber.New()
	app.Delete("/v1/planets/:id", ctr.Delete)

	req := httptest.NewRequest("DELETE", "/v1/planets/"+id, nil)
	resp, _ := app.Test(req)

	planetsDBMock.AssertExpectations(t)

	assert := assert.New(t)
	assert.Equal(404, resp.StatusCode, "they should be equal")
}

func TestIntegrationDontDeletePlanetDueInternalError(t *testing.T) {
	id := "id"
	response := int64(-1)
	planetsDBMock := new(PlanetsDBMock)
	planetsDBMock.On("Delete", mock.Anything, id).Return(response, errors.New(""))

	ctr := PlanetsController{
		PlanetsDomain: &PlanetsDomain{
			PlanetsDB: planetsDBMock,
		},
	}

	app := fiber.New()
	app.Delete("/v1/planets/:id", ctr.Delete)

	req := httptest.NewRequest("DELETE", "/v1/planets/"+id, nil)
	resp, _ := app.Test(req)

	planetsDBMock.AssertExpectations(t)

	assert := assert.New(t)
	assert.Equal(500, resp.StatusCode, "they should be equal")
}
