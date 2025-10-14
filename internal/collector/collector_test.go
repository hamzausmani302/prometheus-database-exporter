package collector

import (
	"testing"

	"github.com/go-gota/gota/dataframe"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/schema"
	"github.com/hamzausmani302/prometheus-database-exporter/pkg/cache"
	"github.com/sirupsen/logrus"
)
func TestCollectorMappingFunction(testing *testing.T) {
	// Create a new collector instance
	logger := logrus.New()
	store := cache.NewLocaltimeCache()
	queries := []*schema.Query{
		&schema.Query{
			Name: "test_query",
			Labels: []schema.Label{
				schema.Label{
					Name:        "test_label1",
					StaticValue: "label_test1",
				},
			},
			Metrics: []schema.Metric{
				schema.Metric{
					Name: "metric_1",
					Type: "gauge",
					Help: "This is metric 1",
					Column: "total_wells",
				},
			},
			
	},
	}
	_collector := NewCollector(logger, &store, queries)
	df := dataframe.LoadRecords([][]string{
		{"well_id", "total_wells"},
		{"well_1", "10"},
	},
	)
	metrics , err := _collector.mapToCollectorMetric(df, *queries[0])
	if err != nil {
		testing.Errorf("Error in mapping to collector metric: %s", err.Error())
	}
	if len(metrics) != 1 {
		testing.Errorf("Error in mapping to collector metric, expected length: %d, got: %d", 1, len(metrics))
	}
	if metrics[0].Name != "test_query_metric_1" {
		testing.Errorf("Error in mapping to collector metric, expected metric name: %s, got: %s", "test_query_metric_1", metrics[0].Name)
	}
}