# Standalone Exporter

There are two binaries built from this repository, and they differ only in whether they respect
the `enableApi` / `enableCollector` flags:

| Binary | `enableApi` / `enableCollector` | Use case |
|---|---|---|
| `cmd/standalone` | Ignored — always runs **both** | Single process, single replica. Simplest way to run the exporter. |
| `cmd/exporter` | Respected | Split deployments — e.g. one replica scraping databases (`enableApi: false`) and separate replicas serving `/metrics` (`enableCollector: false`) behind a shared Redis cache/scheduler storage. |

If you're not intentionally scaling the collector and API separately, use `cmd/standalone`.

## Prerequisites

- [Go 1.24+](https://go.dev/dl/)
- A reachable SQL database to point queries at (PostgreSQL is the currently supported engine)
- Optionally, Redis if you want a shared cache/scheduler store instead of the in-memory defaults

## 1. Clone with submodules

The task scheduler (`pkg/go-scheduler`) is a git submodule — don't skip `--recurse-submodules`.

```bash
git clone --recurse-submodules https://github.com/hamzausmani302/prometheus-database-exporter.git
cd prometheus-database-exporter
# if you already cloned without submodules:
git submodule update --init --recursive
```

## 2. Create a config file

```bash
cp config/config.example.yaml config/config.yaml
```

Edit `config/config.yaml` with real connection details and at least one query — see
[Configure a Data Source](configure-data-source.md) for the full schema, or copy one of the
recipes in [Examples](../examples.md).

## 3. Build

```bash
go build -o bin/exporter ./cmd/standalone
```

(Swap in `./cmd/exporter` instead if you need the split collector/API deployment described above.)

## 4. Run

```bash
CONFIG_FILE_PATH=config/config.yaml ./bin/exporter
```

`CONFIG_FILE_PATH` defaults to `config/config.yaml` if unset. Two other environment variables
override the config file when set: `ENABLE_API` and `ENABLE_COLLECTOR` (both booleans).

On startup the process logs each data source it connects to and each query it registers with the
scheduler. Once it settles, check:

```bash
curl http://localhost:2112/metrics       # your configured query metrics
curl http://localhost:2112/app-metrics   # the exporter's own Go runtime metrics
```

A newly-added query only appears in `/metrics` after its first `queryRefreshTime` tick has
completed and written a result into the cache store — an empty response for a brand-new query for
the first few seconds is expected, not a bug.

## Running with Docker

The published image builds whichever binary is named by the `SERVICE` build arg (`exporter` by
default):

```bash
docker build --build-arg SERVICE=standalone -t prometheus-database-exporter .
docker run -p 2112:2112 \
  -v $(pwd)/config/config.yaml:/app/config/config.yaml \
  prometheus-database-exporter
```

Pre-built images are also published to `hamzausmani021/prometheus-database-exporter` on tagged
releases (see the CI workflows under `.github/workflows/`).

## Prometheus scrape config

```yaml
scrape_configs:
  - job_name: database-exporter
    static_configs:
      - targets: ["localhost:2112"]
```

## Stopping

The process handles `SIGINT`/`SIGTERM`: it stops the scheduler and closes every data source
connection before exiting, so a plain `Ctrl+C` or `docker stop` is enough — no need to kill `-9`.
