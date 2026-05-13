APP_NAME=ela-api

.PHONY: tidy run build test fmt docker-up docker-down

tidy:
go mod tidy

run:
go run ./cmd/api

build:
go build -o bin/$(APP_NAME) ./cmd/api

test:
go test ./...

fmt:
gofmt -w ./cmd ./internal

docker-up:
docker compose up --build -d

docker-down:
docker compose down -v
