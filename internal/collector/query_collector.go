package collector

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/schema"
	"github.com/hamzausmani302/prometheus-database-exporter/pkg/cache"
	"github.com/hamzausmani302/prometheus-database-exporter/pkg/utils"
	"github.com/sirupsen/logrus"
)

// prometheus uses float64 or int
type MetricValue float64
type QueryCollector struct {
	Logger    *logrus.Logger
	DataStore *cache.ICache
	Queries   []*schema.Query
}

func (_collector *QueryCollector) getDataFromStore(key string) (dataframe.DataFrame, error) {
	// fetching data from store
	_collector.Logger.Infof("Getting data for task id = %s", key)
	var bytesData []byte
	if d, err := (*_collector.DataStore).Get(key); err == nil && d != nil {
		_collector.Logger.Debug(d)
		bytesData = d
	} else {
		return dataframe.DataFrame{}, err
	}

	return utils.DataFrameFromCSVBytes(bytesData), nil
}

func (_collector *QueryCollector) assignLabels(cols []string, record []string, query *schema.Query) []CollectorMetricLabel {
	var commonLabels []CollectorMetricLabel = []CollectorMetricLabel{}
	for _, label := range query.Labels {
		// if static value is not provided, assign the column Value
		// if both are empty, log an error
		var commonLabel *CollectorMetricLabel
		if label.IsStaticValue() {
			// assign the static value
			commonLabel = &CollectorMetricLabel{
				Name:  label.Name,
				Value: label.StaticValue,
			}
		} else {
			// assign the dynamic value from column of dataframe
			indx := slices.Index(cols, label.ColumnName)
			if indx != -1 {
				commonLabel = &CollectorMetricLabel{
					Name:  label.Name,
					Value: record[indx],
				}
			} else {
				_collector.Logger.Warnf("label %q: column %q not found in query results", label.Name, label.ColumnName)
			}
		}
		if commonLabel != nil {
			commonLabels = append(commonLabels, *commonLabel)
		}
	}
	return commonLabels
}

// metricFingerprint returns a string key uniquely identifying a metric by its name and label set.
// Used to detect and deduplicate metrics that share the same (name, labels) combination.
func metricFingerprint(name string, labels []CollectorMetricLabel) string {
	parts := make([]string, 0, 1+len(labels))
	parts = append(parts, name)
	for _, l := range labels {
		parts = append(parts, l.Name+"="+l.Value)
	}
	return strings.Join(parts, "|")
}

func (_collector *QueryCollector) mapToCollectorMetric(df dataframe.DataFrame, query schema.Query) ([]CollectorMetric[MetricValue], error) {
	_collector.Logger.Debug("Prometheus mapping to collector")
	cols := df.Names()
	_collector.Logger.Debug(cols)
	records := df.Copy().Records()
	if len(records) <= 1 {
		_collector.Logger.Warn("Result of query is empty")
		return []CollectorMetric[MetricValue]{}, nil
	}

	// Use a map to deduplicate by (metric name, label set). Last value wins.
	seen := make(map[string]int)
	exportMetrics := []CollectorMetric[MetricValue]{}

	for _, metric := range query.Metrics {
		for i := 1; i < len(records); i++ {
			// Assign labels
			labels := _collector.assignLabels(cols, records[i], &query)
			// Map Object to CollectorMetric
			idx := slices.Index(cols, metric.Column)
			if idx == -1 {
				_collector.Logger.Warnf("column %s not present in dataframe", metric.Column)
				continue
			}
			value, err := strconv.ParseFloat(records[i][idx], 64)
			if err != nil {
				_collector.Logger.Errorf("For metric = %s , error converting result %s to float", metric.Name, records[i][idx])
				continue
			}

			metricName := fmt.Sprintf("%s_%s", query.Name, metric.Name)
			fp := metricFingerprint(metricName, labels)

			exportMetric := CollectorMetric[MetricValue]{
				Name:   metricName,
				Labels: labels,
				Value:  MetricValue(value),
				Type:   metric.Type,
				Help:   metric.Help,
			}

			if pos, duplicate := seen[fp]; duplicate {
				_collector.Logger.Debugf("duplicate metric %q with same label set — keeping last value", metricName)
				exportMetrics[pos] = exportMetric
			} else {
				seen[fp] = len(exportMetrics)
				exportMetrics = append(exportMetrics, exportMetric)
			}
		}
	}
	_collector.Logger.Debug(`________________Output Metrics______________`)
	_collector.Logger.Debug(exportMetrics)
	_collector.Logger.Debug(`__________________________________________`)

	return exportMetrics, nil
}

func (_collector *QueryCollector) GetCollectedMetrics() ([]CollectorMetric[MetricValue], error) {
	export_metrics := []CollectorMetric[MetricValue]{}
	for _, query := range _collector.Queries {
		_collector.Logger.Debugf("Query data for hash = %s", query.GetHash())
		df, err := _collector.getDataFromStore(query.GetHash())
		_collector.Logger.Debug("data", df)
		if err != nil {
			_collector.Logger.Error(err)
		}
		metrics, errc := _collector.mapToCollectorMetric(df, *query)
		if errc != nil {
			_collector.Logger.Error(errc)
		}
		export_metrics = append(export_metrics, metrics...)
	}
	return export_metrics, nil
}

func (_collector *QueryCollector) scrapeMetric(metrics []CollectorMetric[MetricValue]) error {
	_collector.Logger.Debug("Prometheus scraping metrics")
	return nil
}

func NewCollector(logger *logrus.Logger, store *cache.ICache, queries []*schema.Query) ICollector[MetricValue] {
	return &QueryCollector{
		Logger:    logger,
		DataStore: store,
		Queries:   queries,
	}
}
