package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlanetsDomainMock struct {
	mock.Mock
}

func (client *PlanetsDomainMock) Get(filter bson.M) ([]Planet, error) {
	args := client.Called(filter)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Planet), args.Error(1)
}

func (client *PlanetsDomainMock) Create(body string) (map[string]string, error) {
	args := client.Called(body)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(map[string]string), args.Error(1)
}

func (client *PlanetsDomainMock) Delete(id string) (int64, error) {
	args := client.Called(id)
	return args.Get(0).(int64), args.Error(1)
}

func TestGetPlanetsSuccessfully(t *testing.T) {
	filter := bson.M{}
	response := []Planet{{Name: "bar", Terrain: "world", Climate: "cold"}}

	planetsDomainMock := new(PlanetsDomainMock)
	planetsDomainMock.On("Get", filter).Return(response, nil)

	ctr := PlanetsController{
		PlanetsDomain: planetsDomainMock,
	}

	app := fiber.New()
	app.Get("/v1/planets", ctr.Index)

	req := httptest.NewRequest("GET", "/v1/planets", nil)
	resp, _ := app.Test(req)

	planetsDomainMock.AssertExpectations(t)

	assert := assert.New(t)
	assert.Equal(200, resp.StatusCode, "they should be equal")
}

func TestDontGetPlanetsDueToNotFound(t *testing.T) {
	filter := bson.M{}

	planetsDomainMock := new(PlanetsDomainMock)
	planetsDomainMock.On("Get", filter).Return([]Planet{}, nil)

	ctr := PlanetsController{
		PlanetsDomain: planetsDomainMock,
	}

	app := fiber.New()
	app.Get("/v1/planets", ctr.Index)

	req := httptest.NewRequest("GET", "/v1/planets", nil)
	resp, _ := app.Test(req)

	planetsDomainMock.AssertExpectations(t)

	assert := assert.New(t)
	assert.Equal(404, resp.StatusCode, "they should be equal")
}

func TestDontGetPlanetsDueToInternalError(t *testing.T) {
	planetsDomainMock := new(PlanetsDomainMock)
	planetsDomainMock.On("Get", mock.Anything).Return(nil, errors.New(""))

	ctr := PlanetsController{
		PlanetsDomain: planetsDomainMock,
	}

	app := fiber.New()
	app.Get("/v1/planets", ctr.Index)

	req := httptest.NewRequest("GET", "/v1/planets", nil)
	resp, _ := app.Test(req)

	planetsDomainMock.AssertExpectations(t)

	assert := assert.New(t)
	assert.Equal(500, resp.StatusCode, "they should be equal")
}

func TestCreatePlanetSuccessfully(t *testing.T) {
	body, _ := json.Marshal(map[string]interface{}{
		"name":    "earth",
		"terrain": "grass",
		"climate": "rainy",
	})

	var response map[string]string
	json.Unmarshal(body, &response)

	planetsDomainMock := new(PlanetsDomainMock)
	planetsDomainMock.On("Create", mock.Anything).Return(response, nil)

	ctr := PlanetsController{
		PlanetsDomain: planetsDomainMock,
	}

	app := fiber.New()
	app.Post("/v1/planets", ctr.Create)

	req := httptest.NewRequest("POST", "/v1/planets", nil)
	resp, _ := app.Test(req)

	planetsDomainMock.AssertExpectations(t)

	assert := assert.New(t)
	assert.Equal(201, resp.StatusCode, "they should be equal")
}

func TestDontCreatePlanetDueToConflict(t *testing.T) {
	writeErrors := []mongo.WriteError{{Code: 11000}}
	planetsDomainMock := new(PlanetsDomainMock)
	planetsDomainMock.On("Create", mock.Anything).Return(nil, mongo.WriteException{WriteErrors: writeErrors})

	ctr := PlanetsController{
		PlanetsDomain: planetsDomainMock,
	}

	app := fiber.New()
	app.Post("/v1/planets", ctr.Create)

	req := httptest.NewRequest("POST", "/v1/planets", nil)
	resp, _ := app.Test(req)

	planetsDomainMock.AssertExpectations(t)

	assert := assert.New(t)
	assert.Equal(409, resp.StatusCode, "they should be equal")
}

func TestDontCreatePlanetDueToBadRequest(t *testing.T) {
	planetsDomainMock := new(PlanetsDomainMock)
	planetsDomainMock.On("Create", mock.Anything).Return(nil, validator.ValidationErrors{})

	ctr := PlanetsController{
		PlanetsDomain: planetsDomainMock,
	}

	app := fiber.New()
	app.Post("/v1/planets", ctr.Create)

	req := httptest.NewRequest("POST", "/v1/planets", nil)
	resp, _ := app.Test(req)

	planetsDomainMock.AssertExpectations(t)

	assert := assert.New(t)
	assert.Equal(400, resp.StatusCode, "they should be equal")
}

func TestDontCreatePlanetDueToInternalError(t *testing.T) {
	planetsDomainMock := new(PlanetsDomainMock)
	planetsDomainMock.On("Create", mock.Anything).Return(nil, errors.New(""))

	ctr := PlanetsController{
		PlanetsDomain: planetsDomainMock,
	}

	app := fiber.New()
	app.Post("/v1/planets", ctr.Create)

	req := httptest.NewRequest("POST", "/v1/planets", nil)
	resp, _ := app.Test(req)

	planetsDomainMock.AssertExpectations(t)

	assert := assert.New(t)
	assert.Equal(500, resp.StatusCode, "they should be equal")
}

func TestDeletePlanetSuccessfully(t *testing.T) {
	id := "124"

	planetsDomainMock := new(PlanetsDomainMock)
	planetsDomainMock.On("Delete", id).Return(int64(1), nil)

	ctr := PlanetsController{
		PlanetsDomain: planetsDomainMock,
	}

	app := fiber.New()
	app.Delete("/v1/planets/:id", ctr.Delete)

	url := fmt.Sprintf("/v1/planets/%s", id)

	req := httptest.NewRequest("DELETE", url, nil)
	resp, _ := app.Test(req)

	planetsDomainMock.AssertExpectations(t)

	assert := assert.New(t)
	assert.Equal(204, resp.StatusCode, "they should be equal")
}

func TestDontDeletePlanetDueToNotFound(t *testing.T) {
	id := "124"

	planetsDomainMock := new(PlanetsDomainMock)
	planetsDomainMock.On("Delete", id).Return(int64(0), errors.New(""))

	ctr := PlanetsController{
		PlanetsDomain: planetsDomainMock,
	}

	app := fiber.New()
	app.Delete("/v1/planets/:id", ctr.Delete)

	url := fmt.Sprintf("/v1/planets/%s", id)

	req := httptest.NewRequest("DELETE", url, nil)
	resp, _ := app.Test(req)

	planetsDomainMock.AssertExpectations(t)

	assert := assert.New(t)
	assert.Equal(404, resp.StatusCode, "they should be equal")
}

func TestDontDeletePlanetDueInternalError(t *testing.T) {
	id := "124"

	planetsDomainMock := new(PlanetsDomainMock)
	planetsDomainMock.On("Delete", id).Return(int64(-1), errors.New(""))

	ctr := PlanetsController{
		PlanetsDomain: planetsDomainMock,
	}

	app := fiber.New()
	app.Delete("/v1/planets/:id", ctr.Delete)

	url := fmt.Sprintf("/v1/planets/%s", id)

	req := httptest.NewRequest("DELETE", url, nil)
	resp, _ := app.Test(req)

	planetsDomainMock.AssertExpectations(t)

	assert := assert.New(t)
	assert.Equal(500, resp.StatusCode, "they should be equal")
}