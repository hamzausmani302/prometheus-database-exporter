.PHONY: build test-unit test-integration test-e2e lint docker-build submodules setup

## Build the exporter binary
build:
	go build -o bin/exporter ./cmd/exporter/

## Run unit tests
test-unit:
	go test ./... -count=1

## Run integration tests (requires Postgres + Redis via docker-compose)
test-integration:
	go test -tags=integration ./test/integration/ -count=1

## Run end-to-end tests (requires services and seeded database)
test-e2e:
	CONFIG_FILE_PATH=config/config.test.yaml go test -tags=e2e ./test/e2e/ -count=1

## Run linter
lint:
	go vet ./...
	golangci-lint run ./...

## Build Docker image
docker-build:
	docker build -t prometheus-database-exporter:local .

## Start local development services (Postgres + Redis)
services-up:
	docker compose -f test/e2e/docker/docker-compose.yml up -d

## Stop local development services
services-down:
	docker compose -f test/e2e/docker/docker-compose.yml down

## Initialize git submodules
submodules:
	git submodule update --init --recursive

## First-time project setup
setup: submodules
	cp -n config/config.example.yaml config/config.yaml || true
	@echo "Setup complete. Edit config/config.yaml with your database details."
