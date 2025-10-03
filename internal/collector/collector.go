package collector

import (
	"github.com/go-gota/gota/dataframe"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/schema"
)

// struct representing actual metric to be exposed to scraping software
type CollectorMetric[T comparable] struct {
	Name   string
	Labels []CollectorMetricLabel
	Help string
	Type string
	Value  T
}
// Struct represening the label for the metric
type CollectorMetricLabel struct {
	Name  string
	Value string
}
/* 
Interface to be followed for the implementation of a different collector
	Output: ([]CollectorMetric[T], error)
*/
type ICollector[T comparable] interface {
	// Get data from store in which the schduler task put it into
	getDataFromStore(key string) (dataframe.DataFrame,error)
	// Converts the data to CollectorMetric[T] form
	mapToCollectorMetric(df dataframe.DataFrame, query schema.Query) ([]CollectorMetric[T], error)
	// the following method will be implemented by subclasses, to prepare data for scraping or to send to external systems
	scrapeMetric( metrics []CollectorMetric[T]) error
	// Collect method use the following method , callable by collector from ourside
	GetCollectedMetrics() ([]CollectorMetric[T], error)
	// Promethsu specofic methods
}

