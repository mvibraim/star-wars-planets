package main

type Planet struct {
	Name    string `json:"name,omitempty"`
	Weather string `json:"weather,omitempty"`
	Terrain string `json:"terrain,omitempty"`
}
