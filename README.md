[![CircleCI](https://circleci.com/gh/mvibraim/star-wars-planets.svg?style=svg)](https://circleci.com/gh/mvibraim/star-wars-planets)

# Star Wars Planets

Create, delete and list planets. If a created planet exists in Star Wars, the movie appearances is saved.

Uses Redis to cache the movie appearances indexed by planet name on application startup.

All data provided by the user through the API is persisted in MongoDB along with movie appearances.

## Technologies used

- Go
- Fiber
- MongoDB
- Redis

## Features

- Create planet
- Get all planets (only `name`, `climate` and `terrain` are shown)
- Get planets filtered by `id` or `name` (only `name`, `climate` and `terrain` are shown)
- Delete planet
- Cache SWAPI planets data as key-value `planetName: moviesAppearances`
- `name` as unique index
- Application logs
- Validate all payload fields as required

## Tests

`go test` run all the tests below

- Planets controller unit test
- Planets cache unit test
- SWAPI client unit test
- Planets domain unit test
- Planets API integration test

## Prerequisites

- [Docker 18.06.0+](https://docs.docker.com/install/)
- [Docker-Compose](https://docs.docker.com/compose/install/)
- [Make](https://www.gnu.org/software/make/)

To check if the prerequisites are installed, just use the following commands:

```bash
docker -v
docker-compose -v
make -v
```

## Usage

1. Run `make`

2. Create a planet

   ```bash
   curl -X POST \
   -d '{
         "name": "Tatooine",
         "terrain": "desert",
         "climate": "arid"
       }' \
   localhost:3000/v1/planets
   ```

   The movie appearances will be logged to console

   ###### Responses
   ```bash
   201 Created     {"id":"5ecbc20616e02079fab641c0"}
   409 Conflict    {"message":"'name' already exists"}
   400 Bad Request {"errors":[{"error":"required","param":"Name"},{"error":"required","param":"Climate"},{"error":"required","param":"Terrain"}]}
   ```

3. Get planets

   ```bash
   curl -X GET localhost:3000/v1/planets
   curl -X GET localhost:3000/v1/planets?id=5ec56d534c835f8177bff423
   curl -X GET localhost:3000/v1/planets?name=Tatooine
   ```

   ###### Responses
   ```bash
   200 OK [{"name":"Tatooine","climate":"arid","terrain":"desert"}, ...}
   200 OK []
   ```

4. Delete planet

   ```bash
   curl -X DELETE localhost:3000/v1/planets/5ec56d534c835f8177bff423
   ```

   ###### Responses
   ```bash
   204 No Content
   404 Not Found
   ```