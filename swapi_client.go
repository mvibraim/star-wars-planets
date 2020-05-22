package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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

func cacheMovieAppearancesByName() {
	fmt.Printf("%s\n", "Caching movie appearances indexed by name")

	allMovieAppearancesIndexedByName := fetchPlanets()
	conn := getRedisConn()

	for _, filmData := range allMovieAppearancesIndexedByName {
		for name, movieAppearances := range filmData {
			setCache(conn, name, movieAppearances)
		}
	}

	fmt.Printf("%s\n", "Cached successfully")

	conn.Close()
}

func fetchPlanets() []map[string]int {
	fmt.Printf("%s\n", "Fetching all planets from SWAPI API and indexing movie appearances by name, as []map[string]int{name: movieAppearances}")

	channel := make(chan []map[string]int)

	go fetchMovieAppearancesIndexedByName(channel)

	allMovieAppearancesIndexedByName := <-channel

	fmt.Printf("%s\n", "Fetched successfully")

	return allMovieAppearancesIndexedByName
}

func fetchMovieAppearancesIndexedByName(channel chan []map[string]int) {
	planetsURL := config.SwapiURL
	var allMovieAppearancesIndexedByName []map[string]int

	for {
		body := fetchSwapiPlanetsBody(planetsURL)

		allMovieAppearancesIndexedByName = IndexMovieAppearancesByName(body, allMovieAppearancesIndexedByName)

		if body.Next == "" {
			break
		} else {
			planetsURL = body.Next
		}
	}

	channel <- allMovieAppearancesIndexedByName
}

func fetchSwapiPlanetsBody(planetsURL string) SwapiPlanetsBody {
	res, _ := http.Get(planetsURL)
	rawBody, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	var body SwapiPlanetsBody
	json.Unmarshal(rawBody, &body)

	return body
}

// IndexMovieAppearancesByName indexes the movie appearances by name
func IndexMovieAppearancesByName(body SwapiPlanetsBody, allMovieAppearancesIndexedByName []map[string]int) []map[string]int {
	for _, planet := range body.Results {
		movieAppearances := len(planet.Films)
		name := planet.Name

		indexedByName := map[string]int{name: movieAppearances}

		allMovieAppearancesIndexedByName = append(allMovieAppearancesIndexedByName, indexedByName)
	}

	return allMovieAppearancesIndexedByName
}
