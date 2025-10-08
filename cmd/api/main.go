package main

import (
	"fmt"
	"net/http"

	"github.com/algorythma/go-scheduler"
	"github.com/hamzausmani302/prometheus-database-exporter/config"
	col "github.com/hamzausmani302/prometheus-database-exporter/internal/collector"
	promcollector "github.com/hamzausmani302/prometheus-database-exporter/internal/collector/prometheus"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/datasource"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/factories"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/queryscheduler"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/schema"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func main() {
	reg := prometheus.NewPedanticRegistry()

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
	// Mapping Query to class Object
	queryCollector := col.MCollector{DataStore: cacheStore, Logger: logger, Queries: queries}
	// create the promethus collector	
	reg.MustRegister(promcollector.PrometheusGoCollector{
		Logger: logrus.New(),
		Collector: &queryCollector,
	})


	http.Handle("/app-metrics", promhttp.Handler())
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	http.ListenAndServe(":2112", nil)	
}