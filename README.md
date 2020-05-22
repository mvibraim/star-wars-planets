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
- Get all planets (only `name`, `climate` and `terrain` is shown)
- Get planets filter by `id` or `name` (only `name`, `climate` and `terrain` is shown)
- Delete planet
- Cache SWAPI planets data as key-value `planetName: moviesAppearances`
- `name` as unique index
- Application logs
  
TODO
- Validate required fields in payload
- List pagination
  
## Tests

Run `go test`

- Main unit test

TODO
- Cache unit test
- Planets controller unit test
- Planets client unit test
- SWAPI client unit test
- Main integration test

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

3. Get planets

   ```bash
   curl -X GET localhost:3000/v1/planets
   curl -X GET localhost:3000/v1/planets?id=5ec56d534c835f8177bff423
   curl -X GET localhost:3000/v1/planets?name=Tatooine
   ```

4. Delete planet

   ```bash
   curl -X DELETE localhost:3000/v1/planets/5ec56d534c835f8177bff423
   ```
