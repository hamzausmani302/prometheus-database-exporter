# Examples

Ready-to-adapt query configs. Each assumes a `dataSourceConfig` entry named `app-db` — swap in your
own data source name. See [Configure a Data Source](getting-started/configure-data-source.md) for
the full field reference these build on.

## Simple gauge

One row, one metric, no labels.

```yaml
queries:
  - name: order_stats
    dataSource: app-db
    query: "SELECT count(*) AS open_orders FROM orders WHERE status = 'open'"
    queryRefreshTime: 30
    queryTimeout: 5
    metrics:
      - name: open_orders
        type: GAUGE
        help: "Number of currently open orders"
        column: open_orders
```

Produces a single series: `order_stats_open_orders`.

## Per-row labels from the result set

When a query returns multiple rows, add a `label` with `columnName` for each row's key — every row
becomes its own labeled series.

```yaml
queries:
  - name: orders_by_region
    dataSource: app-db
    query: "SELECT region, count(*) AS order_count FROM orders WHERE status = 'open' GROUP BY region"
    queryRefreshTime: 30
    queryTimeout: 5
    labels:
      - name: region
        columnName: region
    metrics:
      - name: order_count
        type: GAUGE
        help: "Open orders per region"
        column: order_count
```

Produces `orders_by_region_order_count{region="eu"}`, `orders_by_region_order_count{region="us"}`,
one series per row returned.

## Static labels

Use `staticValue` to tag every series from a query with something that isn't in the result set at
all — commonly the data source or environment name, useful once you're running the same query
against several databases.

```yaml
queries:
  - name: orders_by_region
    dataSource: app-db
    query: "SELECT region, count(*) AS order_count FROM orders WHERE status = 'open' GROUP BY region"
    queryRefreshTime: 30
    queryTimeout: 5
    labels:
      - name: region
        columnName: region
      - name: database
        staticValue: app-db      # same value on every series from this query
    metrics:
      - name: order_count
        type: GAUGE
        help: "Open orders per region"
        column: order_count
```

## Multiple metrics from one query

A single query can populate several metrics at once, as long as every column it needs is in the
result set — cheaper than running the same joins twice.

```yaml
queries:
  - name: job_stats
    dataSource: app-db
    query: |
      SELECT job_name,
             last_duration_seconds,
             last_memory_megabytes
      FROM job_runs
      WHERE finished_at > now() - interval '1 hour'
    queryRefreshTime: 60
    queryTimeout: 10
    labels:
      - name: job_name
        columnName: job_name
    metrics:
      - name: last_duration_seconds
        type: GAUGE
        help: "Duration of the most recent run, in seconds"
        column: last_duration_seconds
      - name: last_memory_megabytes
        type: GAUGE
        help: "Peak memory of the most recent run, in megabytes"
        column: last_memory_megabytes
```

Produces two families of series sharing the same `job_name` label:
`job_stats_last_duration_seconds{job_name="..."}` and
`job_stats_last_memory_megabytes{job_name="..."}`.

## Pipeline example

For "for each row from query A, run a parameterized query against table B" — expressed as an
`extract` stage feeding a `foreach` stage. Here: for every active device, sum its offline gaps.

```yaml
queries:
  - name: device_offline_seconds
    queryRefreshTime: 180
    queryTimeout: 60
    pipeline:
      - stageId: devices
        stageType: extract
        dataSource: app-db
        query: "SELECT device_id FROM devices WHERE is_active"

      - stageId: gaps
        stageType: foreach
        dataSource: app-db
        inputStageIds: [devices]
        query: |
          SELECT sum(gap_seconds) AS device_offline_seconds
          FROM connectivity_gaps
          WHERE device_id = {{device_id}}
    labels:
      - name: device_id
        columnName: device_id
    metrics:
      - name: device_offline_seconds
        type: GAUGE
        help: "Total offline duration in seconds over the tracked window"
        column: device_offline_seconds
```

The `devices` stage runs once; the `gaps` stage runs once *per row* it returned, substituting
`{{device_id}}` each time, then all per-row results are concatenated into the final result set —
see [Architecture § Pipelines & Stages](architecture.md#pipelines-stages) for how stage
dependencies are resolved.

## A note on missing or non-numeric values

If a metric's `column` is `NULL` or otherwise can't parse as a float for a given row, that row is
skipped for that metric (and logged) rather than exported as `0` — so a query that sometimes
returns `NULL` will simply produce fewer series on some scrapes, not a zero-valued one.
