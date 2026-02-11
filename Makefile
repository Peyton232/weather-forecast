APP_NAME=weather-api
PORT=8080

.PHONY: run build test docker-build docker-run tidy clean

run:
	go run ./cmd/server

build:
	go build -o bin/$(APP_NAME) ./cmd/server

test:
	go test ./...

tidy:
	go mod tidy

docker-build:
	docker build -t $(APP_NAME) .

docker-run:
	docker run -p $(PORT):$(PORT) $(APP_NAME)
