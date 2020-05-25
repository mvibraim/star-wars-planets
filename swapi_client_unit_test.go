package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCacheMovieAppearancesByNameSuccessfully(t *testing.T) {
	planetsCacheMock := new(PlanetsCacheMock)
	planetsCacheMock.On("SetCache", "Tatooine", 1).Return(nil)

	swapiData := `{"next":null,"results":[{"name":"Tatooine","films":["film"]}]}`
	body := ioutil.NopCloser(bytes.NewReader([]byte(swapiData)))
	resp := http.Response{Status: "200", Body: body}

	planetsHttpClientMock := new(PlanetsHttpClientMock)
	planetsHttpClientMock.On("Get", mock.Anything).Return(&resp, nil)

	swapiClient := SwapiClient{
		PlanetsCache:      planetsCacheMock,
		PlanetsHttpClient: planetsHttpClientMock,
	}

	err := swapiClient.CacheMovieAppearancesByName()

	planetsCacheMock.AssertExpectations(t)
	planetsHttpClientMock.AssertExpectations(t)

	assert := assert.New(t)

	assert.Equal(nil, err, "they should be equal")
}

func TestCacheMovieAppearancesByNameDueToSetCacheError(t *testing.T) {
	cacheError := errors.New("")

	planetsCacheMock := new(PlanetsCacheMock)
	planetsCacheMock.On("SetCache", "Tatooine", 1).Return(cacheError)

	swapiData := `{"next":null,"results":[{"name":"Tatooine","films":["film"]}]}`
	body := ioutil.NopCloser(bytes.NewReader([]byte(swapiData)))
	resp := http.Response{Status: "200", Body: body}

	planetsHttpClientMock := new(PlanetsHttpClientMock)
	planetsHttpClientMock.On("Get", mock.Anything).Return(&resp, nil)

	swapiClient := SwapiClient{
		PlanetsCache:      planetsCacheMock,
		PlanetsHttpClient: planetsHttpClientMock,
	}

	err := swapiClient.CacheMovieAppearancesByName()

	planetsCacheMock.AssertExpectations(t)
	planetsHttpClientMock.AssertExpectations(t)

	assert := assert.New(t)

	assert.Equal(cacheError, err, "they should be equal")
}

func TestCacheMovieAppearancesByNameDueToFetchError(t *testing.T) {
	fetchError := errors.New("")

	resp := http.Response{}
	planetsHttpClientMock := new(PlanetsHttpClientMock)
	planetsHttpClientMock.On("Get", mock.Anything).Return(&resp, fetchError)

	swapiClient := SwapiClient{
		PlanetsHttpClient: planetsHttpClientMock,
	}

	err := swapiClient.CacheMovieAppearancesByName()

	planetsHttpClientMock.AssertExpectations(t)

	assert := assert.New(t)

	assert.Equal(fetchError, err, "they should be equal")
}
