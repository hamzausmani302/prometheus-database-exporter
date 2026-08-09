package datasource_executor

import (
	"fmt"

	"github.com/go-gota/gota/dataframe"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/datasource"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/schema"
	"github.com/sirupsen/logrus"
)

type IDataSourceExecutor interface {
	Execute(query schema.Query) (dataframe.DataFrame, error)
}

type DataSourceExecutor struct {
	logger        *logrus.Logger
	dataSourceMap map[string]datasource.IDataSource
}

func (d *DataSourceExecutor) Execute(query schema.Query) (dataframe.DataFrame, error) {
	if query.IsPipeline() {
		pipeline := schema.NewPipeline(d.logger)
		pipeline.BuildPipeline(query.Pipeline, d.dataSourceMap)
		return pipeline.RunPipeline()
	}

	ds, ok := d.dataSourceMap[query.DataSource]
	if !ok {
		return dataframe.DataFrame{}, fmt.Errorf("data source %q not found", query.DataSource)
	}
	if err := ds.Connect(); err != nil {
		return dataframe.DataFrame{}, err
	}
	return ds.GetData(datasource.SQLQuery{Query: query.Query}), nil
}

func NewDataSourceExecutor(logger *logrus.Logger, dataSourceMap map[string]datasource.IDataSource) IDataSourceExecutor {
	return &DataSourceExecutor{
		logger:        logger,
		dataSourceMap: dataSourceMap,
	}
}
