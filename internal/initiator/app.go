package initiator

import (
	"net/http"
	"strings"

	"github.com/algorythma/go-scheduler"
	_storage "github.com/algorythma/go-scheduler/storage"
	"github.com/hamzausmani302/prometheus-database-exporter/config"
	col "github.com/hamzausmani302/prometheus-database-exporter/internal/collector"
	promcollector "github.com/hamzausmani302/prometheus-database-exporter/internal/collector/prometheus"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/datasource"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/factories"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/queryscheduler"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/schema"
	"github.com/hamzausmani302/prometheus-database-exporter/pkg/cache"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// Application struct to hold the application state and configurations
type Application struct {
	logger        *logrus.Logger
	store         *cache.ICache
	cfg           *config.ApplicationConfig
	dataSourceMap map[string]datasource.IDataSource
	queries       []*schema.Query
	storage       _storage.TaskStore
	qScheduler    queryscheduler.IQueryScheduler
	registry      *prometheus.Registry
	Done          chan bool
}

func(app *Application) GetConfig() config.ApplicationConfig{
	return *app.cfg
} 
// Setup and initialize the application components by resolving dependencies
func (app *Application) Init() error {
	app.logger = logrus.New()
	app.logger.Debug("Setting up application")
	cfg := config.GetConfig("example", app.logger)
	app.cfg = &cfg
	app.logger.Debug(cfg)

	// initilizing data sources
	app.dataSourceMap = map[string]datasource.IDataSource{}
	for _, dsource := range cfg.DataSource {
		app.dataSourceMap[dsource.Name] = factories.NewDatasourceFactory(app.logger, &cfg).Create(dsource)
	}

	// Mapping Query to class object
	app.queries = schema.LoadMany(app.logger, cfg.Queries, app.dataSourceMap)
	// Initizing cache store
	cacheStore := factories.NewCacheStoreFactory(app.logger, &cfg).Create(cfg.Store)
	app.store = &cacheStore

	// Initializing schduler
	storage, storageErr := factories.NewSchdulerStorageFactory(app.logger, &cfg).Create(cfg.Scheduler)
	app.storage = storage
	if storageErr != nil {
		app.logger.Panic(storageErr)
		return storageErr
	}
	sch := scheduler.New(storage)
	qScheduler := queryscheduler.NewQuerySchduler(app.logger, &cfg, &sch, app.queries, app.store, &app.Done)
	if err := qScheduler.Init(); err != nil {
		app.logger.Panic("cannot initialize the scheduler", err)
		return err
	}
	app.qScheduler = qScheduler
	app.logger.Info("Configuration initilialized successfully")
	return nil

}

// StartCollector starts the metric collection process
func (app *Application) StartCollector() {
	app.logger.Debug("Starting collector")
	if err := app.qScheduler.Start(); err != nil {
		app.logger.Panic("Failed to start collector", err)
	}
}

// registerCollectors registers the metric collectors with Prometheus or additional backends
func (app *Application) registerCollectors() {
	collectorConfig := app.cfg.Collector
	queryCollector := col.MCollector{DataStore: app.store, Logger: app.logger, Queries: app.queries}
	// all collectors will be registered here
	if strings.EqualFold(strings.ToLower(string(collectorConfig.CollectType)), strings.ToLower(string(config.Prometheus))) {
		// register prometheus collector
		app.registry.MustRegister(promcollector.PrometheusGoCollector{
			Logger:    app.logger,
			Collector: &queryCollector,
		})
	}
}

// StartApi starts the API server to expose metrics
func (app *Application) StartApi() {
	// Initialize registry for API
	app.logger.Debug("Starting API")
	app.registry = prometheus.NewPedanticRegistry()
	app.registerCollectors()
	http.Handle("/app-metrics", promhttp.Handler())
	http.Handle("/metrics", promhttp.HandlerFor(app.registry, promhttp.HandlerOpts{}))

	if err := http.ListenAndServe(":2112", nil)	; err != nil {
		panic(err)
	}

}

// CleanUp performs cleanup operations before shutting down the application
func (app *Application) CleanUp() error {
	if err := app.qScheduler.Stop(); err != nil {
		return err
	}
	for _, dsource := range app.dataSourceMap {
		if err := dsource.Close(); err != nil {
			return err
		}
	}
	return nil
}
