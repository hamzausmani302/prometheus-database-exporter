package config

import (
	"os"
	"testing"
)

func TestReadConfig(t *testing.T) {
	manifest :=
		`storeConfig:
  type: local # enum [Local, Redis] can add more implementations later by implemnting the DataStore interface
  metadata:
    connectionDetails: {}
collectorConfig:
  type: prometheus # enum [Prometheus] can add more implementations later by implemnting the ICollector interface
  metadata: {}
serverConfig:
  port: 8080
  numWorkers: 5
dataSourceConfig:
  - name: mysql-datastore
    type: mysql # enum [MySQL, PostgreSQL] can add more implementations later by implemnting the IDataSource interface
    metadata:
      connectionDetails:
        host: localhost
        port: 3306
        username: root
        password: password
        connectionString: ""
  - name: postgres-datastore
    type: postgresql # enum [MySQL, PostgreSQL] can add more implementations later by implemnting the IDataSource interface
    metadata:
      connectionDetails:
        host: localhost
        port: 5432
        username: postgres
        password: password
        connectionString: "sslmode=disable"
queries:
  - name: taxi_rides
    dataSource: mysql-datastore
    query: "select count(*) AS total_active_rides from rides where status = 'active';"
    labels: {}
    query_refresh_time: 30s
    query_timeout: 5s
    metrics:
      - name: database
        type: LABEL
        help: "Database name"
        static_value: "rides_db"
      - name: total_active_rides # metric will be taxi_rides_total_active_rides in prometheus
        type: GAUGE
        help: "Total number of active rides"
        column: total_active_rides`

	var appConfig ApplicationConfig
	appConfig.readConfigData([]byte(manifest))
	if appConfig.Store.StoreType != "local" {
		t.Errorf("Expected: %s  Got: %s ", "local", appConfig.Store.StoreType)
	}
	if len(appConfig.Queries) != 1 {
		t.Errorf("Expected: %d Got: %d", 1, len(appConfig.Queries))
	}

}


func TestEnvInjection(t *testing.T){
  manifest :=
		`storeConfig:
  type: local # enum [Local, Redis] can add more implementations later by implemnting the DataStore interface
  metadata:
    connectionDetails: {}
collectorConfig:
  type: prometheus # enum [Prometheus] can add more implementations later by implemnting the ICollector interface
  metadata: {}
serverConfig:
  port: 8080
  numWorkers: 5
dataSourceConfig:
  - name: mysql-datastore
    type: mysql # enum [MySQL, PostgreSQL] can add more implementations later by implemnting the IDataSource interface
    metadata:
      connectionDetails:
        host: localhost
        port: 3306
        username: root
        password: password
        connectionString: ""
  - name: postgres-datastore
    type: postgresql # enum [MySQL, PostgreSQL] can add more implementations later by implemnting the IDataSource interface
    metadata:
      connectionDetails:
        host: localhost
        port: 5432
        username: postgres
        password: password
        connectionString: "sslmode=disable"
queries:
  - name: taxi_rides
    dataSource: mysql-datastore
    query: "select count(*) AS total_active_rides from rides where status = 'active';"
    labels: {}
    query_refresh_time: 30s
    query_timeout: 5s
    metrics:
      - name: database
        type: LABEL
        help: "Database name"
        static_value: "rides_db"
      - name: total_active_rides # metric will be taxi_rides_total_active_rides in prometheus
        type: GAUGE
        help: "Total number of active rides"
        column: total_active_rides`

	var appConfig ApplicationConfig
	appConfig.readConfigData([]byte(manifest))
  os.Setenv("STORE_TYPE", "TestStore")
  os.Setenv("PORT", "8000")
  ReadEnvVars(&appConfig)
  if appConfig.Port != 8000 {
    t.Errorf("Expected %d , Got %d" , 8000 , appConfig.Port)
  }
  if appConfig.Store.StoreType != "TestStore" {
    t.Errorf("Expected %s , Got %s" , "TestStore" , appConfig.Store.StoreType)
  }

}


