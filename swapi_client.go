package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type SwapiClient struct {
	PlanetsCache      PlanetsCacheHelper
	PlanetsHttpClient HttpClientHelper
}

func CreateSwapiClient() *SwapiClient {
	return &SwapiClient{
		PlanetsCache:      CreatePlanetsCache(),
		PlanetsHttpClient: CreateHttpClient(),
	}
}

// SwapiPlanetsBody represents the useful body fields from SWAPI API response
type SwapiPlanetsBody struct {
	Next    string
	Results []PlanetInfo
}

// PlanetInfo represents the useful planet fields
type PlanetInfo struct {
	Name  string
	Films []string
}

func (sc *SwapiClient) CacheMovieAppearancesByName() error {
	fmt.Printf("%s\n", "Caching movie appearances indexed by name")

	allMovieAppearancesIndexedByName, err := fetchPlanets(sc.PlanetsHttpClient)

	if err != nil {
		fmt.Printf("%s\n", "Can't cache movie appearances due to fetch error")
		return err
	}

	for _, filmData := range allMovieAppearancesIndexedByName {
		for name, movieAppearances := range filmData {
			err := sc.PlanetsCache.SetCache(name, movieAppearances)

			if err != nil {
				fmt.Printf("%s\n", "Can't cache movie appearances due to set cache error")
				return err
			}
		}
	}

	fmt.Printf("%s\n", "Cached successfully")

	return nil
}

func fetchPlanets(client HttpClientHelper) ([]map[string]int, error) {
	fmt.Printf("%s\n", "Fetching all planets from SWAPI API and indexing movie appearances by name, as []map[string]int{name: movieAppearances}")

	allMovieAppearancesIndexedByName, err := fetchMovieAppearancesIndexedByName(client)

	if err != nil {
		fmt.Printf("%s\n", "Can't fetch planets due to error")
		return nil, err
	}

	fmt.Printf("%s\n", "Fetched successfully")

	return allMovieAppearancesIndexedByName, nil
}

func fetchMovieAppearancesIndexedByName(client HttpClientHelper) ([]map[string]int, error) {
	planetsURL := config.SwapiURL
	var allMovieAppearancesIndexedByName []map[string]int

	for {
		body, err := fetchSwapiPlanetsBody(client, planetsURL)

		if err != nil {
			return nil, err
		}

		allMovieAppearancesIndexedByName = IndexMovieAppearancesByName(body, allMovieAppearancesIndexedByName)

		if body.Next == "" {
			break
		} else {
			planetsURL = body.Next
		}
	}

	return allMovieAppearancesIndexedByName, nil
}

func fetchSwapiPlanetsBody(client HttpClientHelper, planetsURL string) (*SwapiPlanetsBody, error) {
	res, err := client.Get(planetsURL)

	if err != nil {
		return nil, err
	}

	rawBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	res.Body.Close()

	var body SwapiPlanetsBody
	json.Unmarshal(rawBody, &body)

	return &body, nil
}

// IndexMovieAppearancesByName indexes the movie appearances by name
func IndexMovieAppearancesByName(body *SwapiPlanetsBody, allMovieAppearancesIndexedByName []map[string]int) []map[string]int {
	for _, planet := range body.Results {
		movieAppearances := len(planet.Films)
		name := planet.Name

		indexedByName := map[string]int{name: movieAppearances}

		allMovieAppearancesIndexedByName = append(allMovieAppearancesIndexedByName, indexedByName)
	}

	return allMovieAppearancesIndexedByName
}
