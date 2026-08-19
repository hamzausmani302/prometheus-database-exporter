# Configure a Data Source

Everything the exporter does is driven by one YAML file (`config/config.yaml` by default, override
with the `CONFIG_FILE_PATH` env var). This page walks through it top to bottom: register a
database as a **data source**, then write a **query** against it that becomes a Prometheus metric.

## 1. Register a data source

```yaml
dataSourceConfig:
  - name: postgres-datastore        # referenced by queries below
    type: SQL                       # currently the only supported type
    metadata:
      connectionDetails:
        connectionString: "postgresql://user:pass@localhost:5432/mydb?sslmode=disable"
        # host / port / username / password are also accepted and used to build
        # a connection string only if connectionString is left empty
        host: localhost
        port: 5432
        username: postgres
        password: postgres
```

You can list as many data sources as you need — each `query` below picks one by `name`. This is
how a single exporter instance can pull metrics from several databases at once.

!!! note "Only PostgreSQL today"
    `type: SQL` currently maps to a Postgres reader (`internal/datasource/postgres_datasource.go`,
    via `lib/pq`). Adding another engine means implementing `IDataSource` and registering it in
    `DatasourceFactory.Create` — see [Contribution](../contributing.md#adding-a-new-data-source).

## 2. Write a query

The simplest query maps one SQL statement to one metric:

```yaml
queries:
  - name: user_count                 # metric prefix -> user_count_total_users
    dataSource: postgres-datastore
    query: "SELECT count(*) AS total_users FROM users WHERE active = true"
    queryRefreshTime: 60              # seconds between runs
    queryTimeout: 5                   # currently informational; not yet enforced on the query itself
    labels: []
    metrics:
      - name: total_users
        type: GAUGE
        help: "Number of active users"
        column: total_users
```

This produces `user_count_total_users` in `/metrics`, refreshed every 60 seconds.

### Field reference

**Query**

| Field | Required | Description |
|---|---|---|
| `name` | yes | Prefix for every metric this query produces (`<name>_<metric.name>`). Also part of the query's dedup hash, so renaming it registers a new scheduled task. |
| `dataSource` | yes, unless using `pipeline` | Name of an entry in `dataSourceConfig`. |
| `query` | yes, unless using `pipeline` | Raw SQL. Mutually exclusive with `pipeline`. |
| `pipeline` | no | A list of stages instead of a single query — see [Pipelines](#pipelines-multi-step-queries) below. |
| `queryRefreshTime` | yes | Seconds between scheduled runs. Also sets the cache TTL (`2x` this value), so a metric silently disappears from `/metrics` if refreshes stop for more than two intervals. |
| `queryTimeout` | yes | Seconds allowed for the query (see note above). |
| `labels` | no | See below. |
| `metrics` | yes | One or more metrics extracted from the same result set — see below. |

**Label** (`labels[]`)

| Field | Description |
|---|---|
| `name` | Prometheus label name. |
| `columnName` | Column in the query result to read the label value from, per row. Accepts `column_name` as a snake_case alias. |
| `staticValue` | A fixed value instead of a column — use this to tag every series from a query with something like the data source name. Takes priority over `columnName` when both are set. |

**Metric** (`metrics[]`)

| Field | Description |
|---|---|
| `name` | Suffix appended to the query name to form the final metric name. |
| `type` | Currently only `GAUGE` is meaningfully exported (every metric is emitted as a Prometheus Gauge regardless of this value). |
| `help` | Help text shown in `/metrics`. |
| `column` | Column in the query result holding the numeric value. Must parse as a float. |

A query can define **multiple metrics** — each is extracted from the *same* result set, once per
row, with the same labels. See [Examples](../examples.md) for a query emitting two metrics
(duration and memory) from one row.

### Rows and labels combine into series

Every row of the result becomes one exported series per metric, labeled from that row's columns.
If two rows resolve to the same metric name + label set, the **last** one wins (rows are
deduplicated in result order) — so make sure your label columns are unique per row unless that's
what you want.

## Pipelines (multi-step queries)

When one SQL statement can't express what you need — most commonly "for each row from query A, run
a parameterized query against table B" — use `pipeline` instead of `query`:

```yaml
queries:
  - name: total_gap_duration_seconds
    queryRefreshTime: 180
    queryTimeout: 60
    pipeline:
      - stageId: wells
        stageType: extract
        dataSource: postgres-datastore
        query: "SELECT well_id FROM wells WHERE is_live"

      - stageId: gaps
        stageType: foreach
        dataSource: postgres-datastore
        inputStageIds: [wells]
        query: "SELECT sum(gap_seconds) AS total_gap_duration_seconds FROM sensor_gaps WHERE well_id = {{well_id}}"
    labels:
      - name: well_id
        columnName: well_id
    metrics:
      - name: total_gap_duration_seconds
        type: GAUGE
        help: "Total sensor gap duration in seconds"
        column: total_gap_duration_seconds
```

`{{well_id}}` is substituted per-row from the `wells` stage's output before each `gaps` query runs.
See [Architecture § Pipelines & Stages](../architecture.md#pipelines-stages) for how stage types
(`extract`, `foreach`, `rename`) and dependency resolution work, and
[Examples](../examples.md#pipeline-example) for the full end-to-end version of this query.

## 3. Choose a cache store and scheduler storage

These aren't per-query, but they determine whether multiple exporter processes can share query
results and scheduling state:

```yaml
storeConfig:
  type: local            # local | redis
  metadata:
    connectionDetails: {}   # host / port / password, when type: redis

schedulerConfig:
  storage: memory         # memory | sqlite | redis | postgres
  metadata:
    connectionDetails: {}   # host/port/password (redis), dbName (sqlite), connectionString (postgres)
```

Use `local` + `memory` for a single standalone process. If you split the collector and API into
separate processes (see [Standalone Exporter](creating-standalone-exporter.md)), both **must**
point at the same Redis (or Postgres, for scheduler storage) instance, or the API will never see
what the collector wrote.

## Validating your config

There's no separate `validate` command — the fastest feedback loop is running the binary with
your config and watching the logs:

```bash
CONFIG_FILE_PATH=config/config.yaml go run ./cmd/standalone
```

A misconfigured `dataSource` reference on a query logs `data source not found` and that query is
skipped; a bad connection string fails on first connect and is logged, but doesn't crash the
process — check the logs for `Error while running task` or `error executing query` if a metric
never shows up.
