package main

import (
	"errors"
	"testing"

	"github.com/rafaeljusto/redigomock"
)

func TestSetCacheSuccessfully(t *testing.T) {
	connMock := redigomock.NewConn()
	connMock.Command("SET", "tatooine", 5).Expect("ok")

	planetsCache := PlanetsCache{
		Conn: connMock,
	}

	err := planetsCache.setCache("Tatooine", 5)

	if err != nil {
		t.Fatal(err)
	}
}

func TestDontSetCacheDueToError(t *testing.T) {
	connMock := redigomock.NewConn()
	connMock.Command("SET", "tatooine", 5).ExpectError(errors.New(""))

	planetsCache := PlanetsCache{
		Conn: connMock,
	}

	err := planetsCache.setCache("Tatooine", 5)

	if err == nil {
		t.Fatal(err)
	}
}

func TestGetCacheSuccessfully(t *testing.T) {
	connMock := redigomock.NewConn()
	connMock.Command("GET", "tatooine").Expect(int64(5))

	planetsCache := PlanetsCache{
		Conn: connMock,
	}

	movieAppearances, err := planetsCache.getCache("Tatooine")

	if err != nil && movieAppearances != 5 {
		t.Fatal(err)
	}
}

func TestDontGetCacheDueToNotFound(t *testing.T) {
	connMock := redigomock.NewConn()
	connMock.Command("GET", "tatooine").Expect(nil)

	planetsCache := PlanetsCache{
		Conn: connMock,
	}

	movieAppearances, err := planetsCache.getCache("Tatooine")

	if err != nil && movieAppearances != -1 {
		t.Fatal(err)
	}
}

func TestDontGetCacheDueToError(t *testing.T) {
	connMock := redigomock.NewConn()
	connMock.Command("GET", "tatooine").ExpectError(errors.New(""))

	planetsCache := PlanetsCache{
		Conn: connMock,
	}

	movieAppearances, err := planetsCache.getCache("Tatooine")

	if err == nil && movieAppearances != -1 {
		t.Fatal(err)
	}
}
