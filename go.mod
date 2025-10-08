module github.com/hamzausmani302/prometheus-database-exporter

go 1.24.7

require (
	github.com/aleksiumish/in-memory-cache v0.0.0-20221207194228-7a96563e9c52
	github.com/go-gota/gota v0.12.0
	github.com/lib/pq v1.10.9
	github.com/mattn/go-sqlite3 v1.14.32
	github.com/prometheus/client_golang v1.23.2
	github.com/redis/go-redis/v9 v9.14.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.11.1
	gopkg.in/yaml.v3 v3.0.1
)
replace github.com/algorythma/go-scheduler => ./pkg/go-scheduler
require (
	github.com/algorythma/go-scheduler v0.0.0
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.66.1 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	gonum.org/v1/gonum v0.9.1 // indirect
	google.golang.org/protobuf v1.36.8 // indirect
)

