<div style="display: flex; align-items: center; gap: 20px;">

<img src="img/favicon.svg" width="50px" height="50px">

<h1>Prometheus Database Exporter</h1>

</div>

A configuration-driven [Prometheus](https://prometheus.io/) exporter that turns SQL queries into
metrics. Instead of writing code for every new metric, you describe a query, the columns that
become labels, and the columns that become metric values in YAML — the exporter schedules the
query, caches its result, and exposes it on `/metrics`.

It is built to scrape metrics *out of* application databases without adding load to them: every
query has its own refresh interval, results are cached between scrapes, and queries can pull from
several databases and be combined into a single metric via aggregation pipelines.

## Why this exporter

- **Query-level refresh intervals** — a cheap `SELECT count(*)` can refresh every 30 seconds while
  an expensive report-style query only runs every 3 hours, all in the same config file.
- **Multi-database queries** — each query targets a named data source, so one exporter instance can
  read from several databases at once.
- **Pipelines for multi-step metrics** — when a single query isn't enough (e.g. "for each active
  row in table A, look something up in table B"), a query can define a small DAG of stages instead.
- **Decoupled scraping** — a scheduler continuously refreshes query results into a cache store
  (in-memory or Redis); the HTTP API only ever reads from that cache, so a slow database never
  blocks or slows down a Prometheus scrape.
- **Pluggable backends** — data sources, cache stores, and scheduler storage are all defined behind
  interfaces, so new database engines or storage backends can be added without touching the core
  pipeline.

## Where to go next

| I want to... | Read |
|---|---|
| Understand how the pieces fit together | [Architecture](architecture.md) |
| Build and run the exporter | [Standalone Exporter](getting-started/creating-standalone-exporter.md) |
| Write my first query config | [Configure a Data Source](getting-started/configure-data-source.md) |
| See ready-to-adapt query configs | [Examples](examples.md) |
| Set up a local dev environment / contribute | [Contribution](contributing.md) |

## Quick start

```bash
git clone --recurse-submodules https://github.com/hamzausmani302/prometheus-database-exporter.git
cd prometheus-database-exporter
cp config/config.example.yaml config/config.yaml
# edit config/config.yaml with your database connection details and queries

go build -o bin/exporter ./cmd/standalone
CONFIG_FILE_PATH=config/config.yaml ./bin/exporter
```

Metrics are then available at `http://localhost:2112/metrics`.

## Project layout

    mkdocs.yml            # Docs site configuration
    config/
        config.example.yaml  # Annotated example configuration
        config.go            # Config schema (source of truth)
    cmd/
        exporter/            # API and/or collector, toggled by config
        standalone/          # Always runs API + collector in one process
    internal/                # Application internals (scheduler, datasource, collector, ...)
    docs/
        index.md             # This page
        ...                  # Other markdown pages, images and other files
