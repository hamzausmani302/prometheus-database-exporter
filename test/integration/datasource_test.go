//go:build integration
// +build integration

package integration_test

import (
	"testing"

	"github.com/hamzausmani302/prometheus-database-exporter/config"
	"github.com/sirupsen/logrus"

	"github.com/hamzausmani302/prometheus-database-exporter/internal/factories"

	"github.com/hamzausmani302/prometheus-database-exporter/internal/datasource"
)

func TestPostgresDataSourceIntegration(t *testing.T) {
	cfg := config.ApplicationConfig{}
	configData := `enableApi: true
enableCollector: false
schedulerConfig:
 storage: "memory" # enum [redis, sqlite, postgres, memory]
 metadata:
  connectionDetails: {}
storeConfig:
 type: "redis" # enum [Local, Redis] can add more implementations later by implemnting the DataStore interface
 metadata: {}
collectorConfig:
 type: prometheus # enum [Prometheus] can add more implementations later by implemnting the ICollector interface
 metadata: {}
serverConfig:
 port: 8080
 numWorkers: 5
dataSourceConfig:
- name: postgres-datastore
  type: SQL # enum [MySQL, PostgreSQL] can add more implementations later by implemnting the IDataSource interface
  metadata:
   connectionDetails:
    host: localhost
    port: 5432
    username: postgres
    password: password
`
	//check connection
	cfg.ReadConfigData([]byte(configData))
	if len(cfg.DataSource) < 1{
		t.Errorf("Expected number of data sources to be > 0 , but received %d", len(cfg.DataSource))
	}
	logger := logrus.New()
	ds := factories.NewDatasourceFactory(logger, &cfg).Create(cfg.DataSource[0])
	if err :=  ds.Connect(); err != nil {
		t.Error(err)
		return
	}
	result := ds.GetData(datasource.SQLQuery{Query: "select 1 AS col, 2 AS col2"})
	if result.Nrow() <= 0{
		t.Errorf("No data received, expected data")
	}	
	record := result.Copy().Records()
	if record[1][0] != "1" || record[1][1] != "2"{
		t.Errorf("got = {%s %s} ,expected = {1, 2}", record[1][0], record[1][1])
	}
}