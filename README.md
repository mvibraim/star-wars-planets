# Star Wars Planets

Create, delete and list planets. If a created planet exists in Star Wars, the number of film appearances is saved as an attribute

## Technologies used

- Go
- MongoDB
- Redis
  
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

1. Run `make` from root project

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
