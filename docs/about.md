# About

**Prometheus Database Exporter** is a configuration-driven exporter that turns arbitrary SQL
queries into Prometheus metrics — no code required to add, change, or remove a metric, just a YAML
entry. It was built to minimize load on application databases (query-level refresh intervals,
cached results between scrapes) while still supporting the messier cases a single `SELECT`
sometimes can't: multiple databases per exporter instance, and multi-step aggregation pipelines
where one query's output parameterizes the next.

See [Architecture](architecture.md) for how it's put together, or
[Standalone Exporter](getting-started/creating-standalone-exporter.md) to run it.

## Status

The core is in active use: configurable data sources, query scheduling, pipelines, a Prometheus
collector, unit/integration/e2e test coverage, and CI for building binaries and Docker images.
Planned next:

- Additional data source engines beyond PostgreSQL
- A CLI utility to print/debug a query's current result without waiting for a scrape
- Visibility into currently-running/scheduled queries

## Links

- **Source:** [github.com/hamzausmani302/prometheus-database-exporter](https://github.com/hamzausmani302/prometheus-database-exporter)
- **Issues / feature requests:** open a GitHub issue or discussion on the repository above
- **Contributing:** see [Contribution](contributing.md)

## License

Released under the [MIT License](https://github.com/hamzausmani302/prometheus-database-exporter/blob/main/LICENSE).
