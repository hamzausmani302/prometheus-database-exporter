package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	col "github.com/hamzausmani302/prometheus-database-exporter/internal/collector"
	promcollector "github.com/hamzausmani302/prometheus-database-exporter/internal/collector/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/hamzausmani302/prometheus-database-exporter/config"

	"github.com/hamzausmani302/prometheus-database-exporter/internal/datasource"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/queryscheduler"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/schema"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/algorythma/go-scheduler"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/factories"
	"github.com/sirupsen/logrus"
)

type T struct{
	logger *logrus.Logger
}
func (t T) TaskWithArgs(message string) {
	t.logger.Println("TaskWithArgs is executed. message:", message)
}

/*
entry point for collector which collects metrics from the database
and puts them in store from where prmetheus can scrape them via API.
*/
func main() {
	reg := prometheus.NewPedanticRegistry()

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    done := make(chan bool, 1)

	logger := logrus.New()
	fmt.Println("Collector started")	
	
	// Read config from file
	cfg := config.GetConfig("example", logger)
	logger.Debug(cfg)

	// create datasources 
	dataSourceMap := map[string]datasource.IDataSource{}
	for _, dsource := range cfg.DataSource{
		dataSourceMap[dsource.Name] = factories.NewDatasourceFactory(logger, &cfg).Create(dsource)
	}
	
	// Mapping Query to class object
	queries := schema.LoadMany(logger, cfg.Queries, dataSourceMap)
	cacheStore := factories.NewCacheStoreFactory(logger, &cfg).Create(cfg.Store)
	// fmt.Println(q)
	// Creating scheduler(Inject dependencies)
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
	go queryscheduler.Start()
	
	// On interrupt, close datasource, and all the connections
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


    go func() {
		// Listens for intended termination and terminate the memory addresses
		logger.Info("triggered executing")
        sig := <-sigs
        logger.Debug(sig)
        done <- true
		// close scheduler
		if err := queryscheduler.Stop(); err != nil{
			logger.Error(err)
		}
		// close storage connections
		if err := storage.Close(); err != nil {
			logger.Error(err)		
		}
		// clear the connection to database
		for _ , dsource := range dataSourceMap {
			if err := dsource.Close(); err != nil {
				logger.Error(err)
			}
		}
		close(sigs)
		close (done)
		logger.Info("closed successfully")
		}()

    logger.Println("awaiting signal")
    <-done
    logger.Println("exiting")
}

