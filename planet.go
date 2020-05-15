package main

type Planet struct {
	_id     string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Weather string `json:"weather,omitempty"`
	Terrain string `json:"terrain,omitempty"`
}
