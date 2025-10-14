package factories

import (
	"github.com/hamzausmani302/prometheus-database-exporter/config"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/datasource"
	"github.com/sirupsen/logrus"
)

/*
Factory class to initiate the configurations and provide abstraction over the creation of
datasource objects
*/
type DatasourceFactory struct {
	logger *logrus.Logger
	cfg    *config.ApplicationConfig
}

func (dsf *DatasourceFactory) Create(dataSourceConfig config.DataSourceConfig) datasource.IDataSource {
	if dataSourceConfig.DataSourceType == "SQL" {
		return datasource.NewPostgresDatasource(dsf.logger, dsf.cfg, dataSourceConfig)
	}
	dsf.logger.Fatalf("Invalid Data source: %s", dataSourceConfig.DataSourceType)
	return nil
}

func NewDatasourceFactory(logger *logrus.Logger, cfg *config.ApplicationConfig) *DatasourceFactory {
	return &DatasourceFactory{logger, cfg}
}
