package main

import (
	"strings"

	"github.com/gomodule/redigo/redis"
)

type PlanetsCacheHelper interface {
	getCache(name string) (int, error)
	setCache(name string, movieAppearances int) error
}

type PlanetsCache struct {
	Conn redis.Conn
}

func CreatePlanetsCache() *PlanetsCache {
	return &PlanetsCache{
		Conn: getRedisConn(),
	}
}

func getRedisConn() redis.Conn {
	return newRedisPool().Get()
}

func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(config.RedisNetwork, config.RedisHost)

			if err != nil {
				panic(err.Error())
			}

			return c, err
		},
	}
}

func (cache *PlanetsCache) setCache(name string, movieAppearances int) error {
	_, err := cache.Conn.Do("SET", strings.ToLower(name), movieAppearances)

	if err != nil {
		return err
	}

	return nil
}

func (cache *PlanetsCache) getCache(name string) (int, error) {
	movieAppearances, err := redis.Int(cache.Conn.Do("GET", strings.ToLower(name)))

	if err == redis.ErrNil {
		return -1, nil
	} else if err != nil {
		return -1, err
	} else {
		return movieAppearances, nil
	}
}
