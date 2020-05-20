package main

import (
	"strings"

	"github.com/gomodule/redigo/redis"
)

func getRedisConn() redis.Conn {
	return newRedisPool().Get()
}

func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")

			if err != nil {
				panic(err.Error())
			}

			return c, err
		},
	}
}

func setCache(conn redis.Conn, name string, filmsCount int) error {
	_, err := conn.Do("SET", strings.ToLower(name), filmsCount)

	if err != nil {
		return err
	}

	return nil
}

func getCache(c redis.Conn, name string) (int, error) {
	filmsCount, err := redis.Int(c.Do("GET", name))

	if err == redis.ErrNil {
		return -1, nil
	} else if err != nil {
		return -1, err
	} else {
		return filmsCount, nil
	}
}
