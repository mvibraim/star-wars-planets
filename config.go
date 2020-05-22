package main

import (
	"fmt"

	"github.com/caarlos0/env"
)

// Config represents all application env vars
type Config struct {
	Port            int    `env:"PORT" envDefault:"3000"`
	MongoDBHost     string `env:"MONGODB_HOST" envDefault:"mongodb://localhost:27017"`
	MongoDBDatabase string `env:"MONGODB_DATABASE" envDefault:"star-wars-planet"`
	RedisHost       string `env:"REDIS_HOST" envDefault:":6379"`
	RedisNetwork    string `env:"REDIS_NETWORK" envDefault:"tcp"`
	SwapiURL        string `env:"SWAPI_URL" envDefault:"https://swapi.dev/api/planets/"`
}

func parseConfig() Config {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	return cfg
}
