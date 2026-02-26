package datasource_executor

import (
	"github.com/go-gota/gota/dataframe"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/datasource"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/schema"
	"github.com/sirupsen/logrus"
)

type IDataSourceExecutor interface {
	Execute(query schema.Query) (dataframe.DataFrame, error)
}

type DataSourceExecutor struct {
	logger 	 *logrus.Logger
	IDataSourceExecutor
	dataSource *datasource.IDataSource
	dataSourceMap map[string]datasource.IDataSource

}

func (d *DataSourceExecutor) Execute(query schema.Query) (dataframe.DataFrame, error) {
	var df dataframe.DataFrame
	if query.IsPipeline() {
		// execute it as a pipeline
		pipeline := schema.NewPipeline(d.logger)
		pipeline.BuildPipeline(query.Pipeline, d.dataSourceMap)
		stage := pipeline.RunPipeline()
		df =  (*stage).GetBaseStage().GetOutput()
	}else{
		// execute it as a query
		ds := d.dataSourceMap[query.DataSource]
		if err := ds.Connect();err != nil {
			return dataframe.DataFrame{}, err
		}
		df = ds.GetData(datasource.SQLQuery{
			Query: query.Query,
		})
	}
	return df, nil
}

func NewDataSourceExecutor(logger *logrus.Logger, dataSourceMap map[string]datasource.IDataSource) IDataSourceExecutor {
	return &DataSourceExecutor{
		logger:        logger,
		dataSourceMap: dataSourceMap,
	}
}