//go:build integration
// +build integration

package integration_test

import (
	"testing"
	"time"

	"github.com/hamzausmani302/prometheus-database-exporter/config"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/factories"
	"github.com/sirupsen/logrus"
)

// Test simple memory caching Store
// Test Redis Cache connectivity
// The known port and host in Github actions

func TestCacheStore(t *testing.T) {
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
    password: postgres
`
	//check connection
	cfg.ReadConfigData([]byte(configData))

	stores := []config.StoreConfig{
		config.StoreConfig{
			StoreType: "redis",
			Metadata: config.StoreConfigMetadataConfig{
				ConnectionDetails: map[string]string{
					"host": "localhost",
					"port": "6379",
				},
			},
		},
		config.StoreConfig{
			StoreType: "local",
		},
	}
	logger := logrus.New()
	for _, store := range stores {
		cacheStore, err := factories.NewCacheStoreFactory(logger, &cfg).Create(store)
		if err != nil {
			t.Errorf("unable to create cache store: %v", err)
			return
		}
		expected := "testValue123"
		key := "testKey"
		if cacheStore == nil {
			t.Errorf("unable to establish connection to Redis")
			return
		}
		if err := cacheStore.Set(key, []byte(expected), 5); err != nil {
			t.Error("Cannot set the value to redis", err)
		}
		var data string
		// try to retrieve instantly, It should be received
		if value, err := cacheStore.Get(key); err != nil {
			t.Error(err)
		} else {
			data = string(value)
		}
		if data != expected {
			t.Errorf("Expected = %s, Got = %s", expected, data)
		}
		// the value should be expired by now
		timer1 := time.NewTimer(10 * time.Second)
		<-timer1.C
		if value, _ := cacheStore.Get(key); value != nil {
			t.Errorf("The value should be nil: Expected: %s, Got: %s", "nil", value)
		}
	}

}
