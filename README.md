# Database Prometheus Exporter

A flexible and extensible exporter for scraping Prometheus metrics directly from SQL databases. This project is designed to minimize database overhead, handle complex metric scenarios, and support aggregation pipelines across multiple data sources.

## Features

- **Query-Level Refresh Intervals:**  
  Each metric query can be refreshed at its own interval (e.g., complex queries every 3 hours, simple queries more frequently), reducing unnecessary load on your database.

- **Null and Static Metric Handling:**  
  Gracefully handles `NULL` values and static metrics returned from the database, ensuring accurate Prometheus metric output.

- **Multi-Database Support:**  
  Configure each query to target a different database, enabling metrics collection from heterogeneous environments.

- **Aggregation Pipelines:**  
  Define pipelines to aggregate data from multiple sources and output a single metric, simplifying complex metric calculations.

## Getting Started

### Prerequisites

- Node.js (or your project's runtime)
- Supported SQL databases (e.g., PostgreSQL, MySQL, MSSQL)

### Installation

```bash
git clone https://github.com/yourusername/database-prometheus-exporter.git
cd database-prometheus-exporter
npm install
```

### Configuration

Define your metrics and queries in a configuration file (e.g., `config.yaml`):

```yaml
metrics:
  - name: user_count
    query: SELECT COUNT(*) FROM users
    refresh_interval: 60 # seconds
    database: postgres_main

  - name: sales_total
    query: SELECT SUM(amount) FROM sales
    refresh_interval: 10800 # 3 hours
    database: mysql_sales

pipelines:
  - name: total_activity
    sources:
      - user_count
      - sales_total
    aggregation: sum
```

### Running the Exporter

```bash
npm start
```

Metrics will be exposed at `/metrics` endpoint for Prometheus to scrape.

## Handling Null and Static Values

- Null values are converted to `0` or omitted based on configuration.
- Static metrics (unchanging values) are cached and refreshed only at the specified interval.

## Aggregation Pipelines

Pipelines allow you to combine metrics from different sources using simple aggregation functions (e.g., sum, average).

## Contributing

We welcome contributions! Please open issues or submit pull requests for new features, bug fixes, or documentation improvements.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/my-feature`)
3. Commit your changes
4. Push to the branch (`git push origin feature/my-feature`)
5. Open a pull request

## License

This project is open source under the [MIT License](LICENSE).

## Contact

For questions or support, open an issue or reach out via GitHub Discussions.
