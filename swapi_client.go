package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SwapiPlanetsBody struct {
	Next    string
	Results []PlanetInfo
}

type PlanetInfo struct {
	Name  string
	Films []string
}

const url = "https://swapi.dev/api/planets/"

func cacheFilmsCountByName() {
	fmt.Printf("%s\n", "Caching films count indexed by name")

	allFilmsCountIndexedByName := fetchPlanets()
	conn := getRedisConn()

	for _, filmData := range allFilmsCountIndexedByName {
		for name, filmsCount := range filmData {
			setCache(conn, name, filmsCount)
		}
	}

	fmt.Printf("%s\n", "Cached successfully")

	conn.Close()
}

func fetchPlanets() []map[string]int {
	fmt.Printf("%s\n", "Fetching all planets from SWAPI API and indexing filmsCount by name, as []map[string]int{name: filmsCount}")

	channel := make(chan []map[string]int)

	go fetchFilmsCountIndexedByName(channel)

	allFilmsCountIndexedByName := <-channel

	fmt.Printf("%s\n", "Fetched successfully")

	return allFilmsCountIndexedByName
}

func fetchFilmsCountIndexedByName(channel chan []map[string]int) {
	planetsURL := url
	var allFilmsCountIndexedByName []map[string]int

	for {
		body := fetchSwapiPlanetsBody(planetsURL)

		allFilmsCountIndexedByName = IndexFilmsCountByName(body, allFilmsCountIndexedByName)

		if body.Next == "" {
			break
		} else {
			planetsURL = body.Next
		}
	}

	channel <- allFilmsCountIndexedByName
}

func fetchSwapiPlanetsBody(planetsURL string) SwapiPlanetsBody {
	res, _ := http.Get(planetsURL)
	rawBody, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	var body SwapiPlanetsBody
	json.Unmarshal(rawBody, &body)

	return body
}

func IndexFilmsCountByName(body SwapiPlanetsBody, allFilmsCountIndexedByName []map[string]int) []map[string]int {
	for _, planet := range body.Results {
		filmsCount := len(planet.Films)
		name := planet.Name

		indexedByName := map[string]int{name: filmsCount}

		allFilmsCountIndexedByName = append(allFilmsCountIndexedByName, indexedByName)
	}

	return allFilmsCountIndexedByName
}
