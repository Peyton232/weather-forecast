# weather-forecast

Given a set of coordinates, return the forecast for that location

## Make Commands

Run locally:
make run

Build binary:
make build

Run tests:
make test

Build Docker image:
make docker-build

Run Docker container:
make docker-run

### Example

curl "http://localhost:8080/forecast?lat=39.7456&lon=-97.0892"

## Postman Collection

A Postman collection is included in `/query`.

Import:

- `weather-api.postman_collection.json`
- `local.postman_environment.json`

Set environment to `Local Weather API` and run the requests.
