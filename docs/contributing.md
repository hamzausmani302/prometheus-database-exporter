# Contributing

Thanks for your interest in contributing! This page covers getting the project running locally and
submitting a quality pull request. The canonical copy of this guide lives at
[`CONTRIBUTING.md`](https://github.com/hamzausmani302/prometheus-database-exporter/blob/main/CONTRIBUTING.md)
in the repository root.

## Prerequisites

- [Go 1.24+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/) and Docker Compose (for integration/e2e dependencies)
- [golangci-lint](https://golangci-lint.run/usage/install/)

## Local setup

```bash
# The scheduler is a git submodule — required
git clone --recurse-submodules https://github.com/hamzausmani302/prometheus-database-exporter.git
cd prometheus-database-exporter

cp config/config.example.yaml config/config.yaml
# edit config/config.yaml with your database connection details

go build -o bin/exporter ./cmd/exporter/
CONFIG_FILE_PATH=config/config.yaml ./bin/exporter
```

## Running tests

=== "Unit"

    No external services required:

    ```bash
    go test ./...
    ```

=== "Integration"

    Requires Postgres and Redis running locally:

    ```bash
    docker compose -f test/e2e/docker/docker-compose.yml up -d
    go test -tags=integration ./test/integration/
    ```

=== "E2E"

    With the same Docker Compose services running, initialize the test database first:

    ```bash
    PGPASSWORD=password psql -U postgres -h localhost -p 5433 -c "CREATE DATABASE exporter"
    PGPASSWORD=password psql -U postgres -h localhost -p 5433 -d exporter \
      -f test/sql/datasource-test/setup-dummy-data-postgres.sql

    CONFIG_FILE_PATH=config/config.test.yaml go test -tags=e2e ./test/e2e/
    ```

Tear down with:

```bash
docker compose -f test/e2e/docker/docker-compose.yml down
```

## Code style

**Error handling** — library code (`internal/`, `pkg/`) must *return* errors to callers. Don't
`panic()` or `log.Fatal()` outside of `main()`:

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

**Logging** — use the `logrus` logger passed via the struct, never `fmt.Println` or the standard
`log` package in production code paths.

## Adding a new data source

1. Implement the `IDataSource` interface in a new file under `internal/datasource/`.
2. Register the new type constant in `config/config.go`.
3. Add a case to `DatasourceFactory.Create()` in `internal/factories/datasource_factory.go`.
4. Add an integration test under `test/integration/`.

## Adding a new store backend

1. Implement the `ICache` interface in `pkg/cache/`.
2. Add a case to `CacheStoreFactory.Create()` in `internal/factories/store_factory.go`.

## Before opening a PR

1. `go test ./...` passes
2. `go vet ./...` is clean
3. `golangci-lint run` passes
4. No stray `fmt.Println` debug output left in production code
5. No real credentials or connection strings committed

Fill out the pull request template completely, including the testing steps you followed.
