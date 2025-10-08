package collector

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/go-gota/gota/dataframe"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/schema"
	"github.com/hamzausmani302/prometheus-database-exporter/pkg/cache"
	"github.com/hamzausmani302/prometheus-database-exporter/pkg/utils"
	"github.com/sirupsen/logrus"
)

// prometheus uses float64 or int
type MMetricResultType float64;
type MCollector struct {
	Logger *logrus.Logger
	DataStore cache.ICache
	Queries []*schema.Query
	
}

func (_collector *MCollector) getDataFromStore(key string) (dataframe.DataFrame, error) {
	// fetching data from store
	_collector.Logger.Infof("Getting data for task id = %s", key)
	var bytesData []byte;
	if d, err := _collector.DataStore.Get(key); err == nil && d != nil {
		_collector.Logger.Debug(d)
		bytesData = d
	}else{
		return dataframe.DataFrame{}, err
	}

	return utils.DataFrameFromCSVBytes(bytesData), nil
}

func (_collector *MCollector) assignLabels(cols []string, record []string, query *schema.Query) []CollectorMetricLabel{
	var commonLabels []CollectorMetricLabel = []CollectorMetricLabel{};
	for _, label := range query.Labels{
			// if static value is not provided, assign the column Value
			// if both are empty, log an error
			var commonLabel *CollectorMetricLabel
			if label.IsStaticValue(){
				// assign the static value
				commonLabel  = &CollectorMetricLabel{
					Name: label.Name,
					Value: label.StaticValue,
				} 
			}else{
				// assign the dynamic value from column of dataframe
				indx := slices.Index(cols, label.ColumnName)
				if indx != -1{ 
					commonLabel = &CollectorMetricLabel{
						Name: label.Name,
						Value: record[indx],
					}					
				}

			}
			if commonLabel != nil {
				commonLabels = append(commonLabels, *commonLabel)
			}
	}
	return commonLabels
} 

func (_collector *MCollector) mapToCollectorMetric(df dataframe.DataFrame, query schema.Query) ([]CollectorMetric[MMetricResultType], error) {
	_collector.Logger.Debug("Prmehteus mapping to collector")
	cols := df.Names()
	_collector.Logger.Debug(cols)
	records := df.Copy().Records()
	if len(records ) <= 1 {
		_collector.Logger.Warn("Result of query is empty")
		return []CollectorMetric[MMetricResultType]{}, nil
	}
	
	exportMetrics := []CollectorMetric[MMetricResultType]{}
	for _, metric := range query.Metrics {		
		for i:=1; i<len(records); i++ {
			// Assign labels
			labels := _collector.assignLabels(cols,records[i], &query)
			// Map Object to CollectorMetric
			idx := slices.Index(cols, metric.Column)
			if idx != -1 {
				value, err := strconv.ParseFloat(records[i][idx], 64) 
				if err != nil {
					_collector.Logger.Error("For metric = %s , error converting result %s to float", metric.Name, records[i][idx])
				}else{
					exportMetric := CollectorMetric[MMetricResultType]{
						Name: fmt.Sprintf("%s_%s", query.Name, metric.Name),
						Labels: labels,
						Value: MMetricResultType(value),
						Type: metric.Type,
						Help: metric.Help,

					}
					exportMetrics = append(exportMetrics, exportMetric)
				}
			}else{
				_collector.Logger.Warn("column %s not present in dataframe", metric.Column)
			}
 
		}
	}
	_collector.Logger.Debugf(`________________Output Metrics______________
	\n%s\n
	__________________________________________`, exportMetrics)
	
	return exportMetrics, nil
}

func (_collector *MCollector ) GetCollectedMetrics() ([]CollectorMetric[MMetricResultType], error) {
	export_metrics := []CollectorMetric[MMetricResultType]{}
	for _, query := range _collector.Queries{
		_collector.Logger.Debugf("Query data for hash = %s",query.GetHash() )
		df,err := _collector.getDataFromStore(query.GetHash())
		_collector.Logger.Debug("data", df)
		if err != nil {
			_collector.Logger.Error(err)
		}
		metrics,errc := _collector.mapToCollectorMetric(df, *query)
		if errc != nil{
			_collector.Logger.Error(errc)
		}
		export_metrics = append(export_metrics, metrics...)
	}
	return export_metrics, nil
}

func (_collector *MCollector) scrapeMetric(metrics []CollectorMetric[MMetricResultType]) error {
	fmt.Println("Promethus scraping Metric")
	
	
	return nil
}


func  Collect[T comparable](_collector ICollector[MMetricResultType], queries []*schema.Query) ([]CollectorMetric[MMetricResultType], error) {
	result, err := _collector.GetCollectedMetrics(); 
	if err != nil {
		return []CollectorMetric[MMetricResultType]{}, err
	}
	return result, nil
}

