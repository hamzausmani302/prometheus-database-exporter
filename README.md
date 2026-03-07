# Database Prometheus Exporter

A flexible and extensible exporter for scraping Prometheus metrics directly from SQL databases. Designed to minimize database overhead, handle complex metric scenarios, and support aggregation pipelines across multiple data sources.

## Features

- **Query-Level Refresh Intervals:** Each metric query refreshes at its own interval (e.g., complex queries every 3 hours, simple ones more frequently), reducing unnecessary database load.
- **Null and Static Metric Handling:** Gracefully handles `NULL` values and static metrics, ensuring accurate Prometheus output.
- **Multi-Database Support:** Configure each query to target a different database, enabling metrics collection from heterogeneous environments.
- **Aggregation Pipelines:** Define DAG-based pipelines to aggregate data from multiple sources and output a single metric.

## Getting Started

### Prerequisites

- [Go 1.24+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/) (for running dependencies locally)
- A supported SQL database (currently PostgreSQL)

### Installation

```bash
git clone https://github.com/yourusername/database-prometheus-exporter.git
cd database-prometheus-exporter
git submodule update --init --recursive
```

### Configuration

Copy the example config and fill in your database details:

```bash
cp config/config.example.yaml config/config.yaml
```

Edit `config/config.yaml`. The key sections are:

```yaml
dataSourceConfig:
  - name: my-postgres
    type: SQL
    metadata:
      connectionDetails:
        host: localhost
        port: 5432
        username: postgres
        password: your-password
        connectionString: "postgresql://postgres:your-password@localhost:5432/mydb?sslmode=disable"

queries:
  - name: active_users
    dataSource: my-postgres
    query: "SELECT status, COUNT(*) AS user_count FROM users GROUP BY status"
    queryRefreshTime: 60   # seconds between refreshes
    queryTimeout: 10
    labels:
      - name: status
        columnName: status
    metrics:
      - name: user_count
        type: GAUGE
        help: "Number of users per status"
        column: user_count
```

### Building

```bash
go build -o bin/exporter ./cmd/exporter/
```

### Running

```bash
CONFIG_FILE_PATH=config/config.yaml ./bin/exporter
```

Metrics are exposed at `http://localhost:2112/metrics` for Prometheus to scrape.

## Configuration Reference

| Field | Description |
|---|---|
| `enableCollector` | Start the query scheduler that fetches metrics from the database |
| `enableApi` | Start the HTTP server exposing `/metrics` |
| `schedulerConfig.storage` | Task scheduler backend: `memory`, `redis`, `sqlite` |
| `storeConfig.type` | Metrics cache: `local` or `redis` |
| `collectorConfig.type` | Collector implementation: `prometheus` |
| `dataSourceConfig` | List of database connections |
| `queries` | List of SQL queries with label/metric mappings |

## Running Tests

```bash
# Unit tests
go test ./...

# Integration tests (requires Postgres + Redis)
docker compose -f test/e2e/docker/docker-compose.yml up -d
go test -tags=integration ./test/integration/

# E2E tests
go test -tags=e2e ./test/e2e/
```

## Linting

```bash
golangci-lint run
```

## Contributing

We welcome contributions! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for setup instructions, coding standards, and the PR process.

## License

This project is open source under the [MIT License](LICENSE).

## Contact

For questions or support, open an issue or reach out via GitHub Discussions.
