package collector

import (
	"fmt"

	"github.com/go-gota/gota/dataframe"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/schema"
	"github.com/hamzausmani302/prometheus-database-exporter/pkg/cache"
	"github.com/hamzausmani302/prometheus-database-exporter/pkg/utils"
	"github.com/sirupsen/logrus"
)

// prometheus uses float64 or int
type PromethusMetricResultType float64;
type PrometheusCollector struct {
	Logger *logrus.Logger
	DataStore cache.ICache
	queries []*schema.Query
}

func (_collector *PrometheusCollector) getDataFromStore(key string) (dataframe.DataFrame, error) {
	// fetching data from store
	_collector.Logger.Debugf("Getting data for task id = %s", key)
	var bytesData []byte;
	if d, err := _collector.DataStore.Get(key); err == nil {
		bytesData = d
	}else{
		return dataframe.DataFrame{}, err
	}
	return utils.DataFrameFromCSVBytes(bytesData), nil
}

func (_collector *PrometheusCollector) mapToCollectorMetric(df dataframe.DataFrame, query schema.Query) ([]CollectorMetric[PromethusMetricResultType], error) {
	fmt.Println("Prmehtus mapping to collector")
	
	return []CollectorMetric[PromethusMetricResultType]{}, nil
}


func (_collector *PrometheusCollector) scrapeMetric(metrics []CollectorMetric[PromethusMetricResultType]) error {
	fmt.Println("Promethus scraping Metric")
	
	return nil
}


func  Collect[T comparable](_collector ICollector[PromethusMetricResultType], queries []*schema.Query) error {
	for _, query := range queries{
		df,err := _collector.getDataFromStore(query.GetHash())
		if err != nil {
			// return err
		}
		metrics,errc := _collector.mapToCollectorMetric(df, *query)
		if errc != nil{
			// return errc
		}
		err = _collector.scrapeMetric(metrics)
	}
	return nil
}

