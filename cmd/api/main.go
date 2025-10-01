package main

import (
	"fmt"

	"github.com/hamzausmani302/prometheus-database-exporter/config"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/collector"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/datasource"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/factories"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/queryscheduler"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/schema"
	"github.com/hamzausmani302/prometheus-database-exporter/pkg/go-scheduler"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	fmt.Println("Collector started")
    done := make(chan bool, 1)

	// Read config from file
	cfg := config.GetConfig("example", logger)
	logger.Debug(cfg)
	cacheStore := factories.NewCacheStoreFactory(logger, &cfg).Create(cfg.Store)

	// create datasources
	dataSourceMap := map[string]datasource.IDataSource{}
	for _, dsource := range cfg.DataSource{
		dataSourceMap[dsource.Name] = factories.NewDatasourceFactory(logger, &cfg).Create(dsource)
	}
	// fetching queries
	queries := schema.LoadMany(logger, cfg.Queries, dataSourceMap)
	//Creating schduler for generating IDS
	storage, storageErr := factories.NewSchdulerStorageFactory(logger, &cfg).Create(cfg.Scheduler)
	if storageErr != nil {
		logger.Panic(storageErr)
		return;
	}
	sch := scheduler.New(storage)
	queryscheduler := queryscheduler.NewQuerySchduler(logger, &cfg, &sch, queries, &cacheStore,  &done )
	if err := queryscheduler.Init(); err != nil {
		logger.Panic("cannot initialize the scheduler", err)
		return;
	}
	
	// Mapping Query to class object
	promCollector := collector.PrometheusCollector{DataStore: cacheStore, Logger: logger}
	if err := collector.Collect[string](&promCollector, queries); err != nil {
		fmt.Println("error", err)
	}
}