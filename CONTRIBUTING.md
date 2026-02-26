# Contributing to Database Prometheus Exporter

Thank you for your interest in contributing! This guide covers everything you need to get the project running locally and submit a quality pull request.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Local Setup](#local-setup)
- [Running Tests](#running-tests)
- [Code Style](#code-style)
- [Submitting a PR](#submitting-a-pr)
- [Adding a New Database](#adding-a-new-database)

---

## Prerequisites

- [Go 1.24+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/) and Docker Compose (for running integration/e2e dependencies)
- [golangci-lint](https://golangci-lint.run/usage/install/) (for linting)

---

## Local Setup

```bash
# 1. Clone the repo with submodules (required — the scheduler is a git submodule)
git clone https://github.com/yourusername/database-prometheus-exporter.git
cd database-prometheus-exporter
git submodule update --init --recursive

# 2. Copy the example config
cp config/config.example.yaml config/config.yaml

# 3. Edit config/config.yaml with your database connection details

# 4. Build the binary
go build -o bin/exporter ./cmd/exporter/

# 5. Run the exporter
CONFIG_FILE_PATH=config/config.yaml ./bin/exporter
```

---

## Running Tests

### Unit Tests

No external services required:

```bash
go test ./...
```

### Integration Tests

Requires Postgres and Redis running locally. Start them with Docker Compose:

```bash
docker compose -f test/e2e/docker/docker-compose.yml up -d
```

Then run:

```bash
go test -tags=integration ./test/integration/
```

### E2E Tests

With the same Docker Compose services running, initialize the test database:

```bash
# Assuming docker-compose is up
PGPASSWORD=password psql -U postgres -h localhost -p 5433 -c "CREATE DATABASE exporter"
PGPASSWORD=password psql -U postgres -h localhost -p 5433 -d exporter -f test/sql/datasource-test/setup-dummy-data-postgres.sql
```

Then run:

```bash
CONFIG_FILE_PATH=config/config.test.yaml go test -tags=e2e ./test/e2e/
```

### Teardown

```bash
docker compose -f test/e2e/docker/docker-compose.yml down
```

---

## Code Style

### Error Handling

Library code (`internal/`, `pkg/`) must **return errors** to callers. Do not use `panic()` or `log.Fatal()` outside of `main()`.

```go
// Good
func (p *PostgresDataSource) Connect() error {
    if _, err := p.reader.Connect(); err != nil {
        return fmt.Errorf("connecting to postgres: %w", err)
    }
    return nil
}

// Bad — kills the process from library code
func (p *PostgresDataSource) Connect() error {
    if _, err := p.reader.Connect(); err != nil {
        panic(err)
    }
    return nil
}
```

### Logging

Use the `logrus` logger passed via the struct — never use `fmt.Println` or the standard `log` package in production code paths.

```go
// Good
p.Logger.Debugf("evaluating stage %s", stage.StageId)

// Bad
fmt.Println("evaluating stage", stage.StageId)
```

### Adding a New Data Source

1. Create a new file in `internal/datasource/` implementing the `IDataSource` interface.
2. Register the new type constant in `config/config.go`.
3. Add a case to `DatasourceFactory.Create()` in `internal/factories/datasource_factory.go`.
4. Add an integration test in `test/integration/`.

### Adding a New Store Backend

1. Implement the `ICache` interface in `pkg/cache/`.
2. Add a case to `CacheStoreFactory.Create()` in `internal/factories/store_factory.go`.

---

## Submitting a PR

Before opening a pull request:

1. **All unit tests pass:** `go test ./...`
2. **No vet errors:** `go vet ./...`
3. **Lint passes:** `golangci-lint run`
4. **No debug output:** ensure no `fmt.Println` calls remain in production code
5. **Credentials:** ensure no real passwords or connection strings are committed

Please fill out the pull request template completely, including the testing steps you followed.
