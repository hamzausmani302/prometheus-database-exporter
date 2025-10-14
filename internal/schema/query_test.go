package schema

import (
	"testing"

	"github.com/hamzausmani302/prometheus-database-exporter/internal/datasource"
	"github.com/sirupsen/logrus"
)

// create query object
// validate the query by loading it
func TestQuery(t *testing.T) {
	var queryYAML map[string]interface{} = map[string]interface{}{
		
			"name":            "taxi_rides",
			"dataSource":      "postgres-datastore",
			"query":           "select 'test1' AS well_id, 10 AS total_wells from t_well where to_load = false",
			"queryRefreshTime": 30,
			"queryTimeout":    10,
			"labels": []map[string]interface{}{
				map[string]interface{}{
					"name": "database",
					"staticValue": "prod",
				},
				map[string]interface{}{
					"name":      "well_id",
					"columnName": "well_id",
				},
			},
	}
	sources := map[string]datasource.IDataSource{
		"postgres-datastore": &datasource.PostgresDataSource{},
	}
	var q *Query = &Query{}
	err := q.Load(&logrus.Logger{}, queryYAML, sources);
	if err != nil {
		t.Errorf("Error while loading query : %s", err.Error())
	}
	
	if q.Name != "taxi_rides" {
		t.Errorf("Error in parsing query name expected : %s, got: %s", "taxi_rides", q.Name)
	}
	if q.QueryRefreshTime != 30 {
		t.Errorf("Error in parsing query refresh time expected : %d, got: %d", 30, q.QueryRefreshTime)
	}
	if len(q.Labels) != 2 {
		t.Errorf("Error in parsing query labels expected : %d, got: %d", 2, len(q.Labels))
	}
	if q.Labels[0].Name != "database" || !q.Labels[0].IsStaticValue() || q.Labels[0].StaticValue != "prod" {
		t.Errorf("Error in parsing query label 1 expected : %s, got: %s", "database:prod", q.Labels[0])
	}	
}